package servitor

import (
  "sync"
  "time"
  "fmt"

  "Guenhwyvar/lib/mlog"
  "Guenhwyvar/bringer"
)

type CircularBuffer struct {
  messages []mlog.Mlog
  size int
  index int
  count int
  mu sync.Mutex
}

func NewCircularBuffer(size int) *CircularBuffer{
  return &CircularBuffer{
    messages: make([]mlog.Mlog, size),
    size: size,
  }
}

func (cb *CircularBuffer) Add(msg mlog.Mlog) {
  cb.mu.Lock()
  defer cb.mu.Unlock()

  cb.messages[cb.index] = msg
  cb.index = (cb.index +1) % cb.size
  if cb.count < cb.size {
    cb.count++
  }
}

func (cb *CircularBuffer) GetLastN(n int) []mlog.Mlog {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if n > cb.count {
        n = cb.count
    }
    result := make([]mlog.Mlog, n)

    for i := 0; i < n; i++ {
        idx := (cb.index - 1 - i + cb.size) % cb.size
        result[n-1-i] = cb.messages[idx]
    }

    return result
}

type MemoryManagerServ struct {
    bringer      bringer.Mesmerizer
    buffer       *CircularBuffer
    groupTZ      *time.Location
    serverTZ     *time.Location
    currentDate  time.Time
    dailyMessages []mlog.Mlog
    mu           sync.Mutex
}

func NewMemoryManagerServ(bringer bringer.Mesmerizer) *MemoryManagerServ {
    groupTZ, _ := time.LoadLocation("Asia/Bangkok") // GMT+7
    serverTZ, _ := time.LoadLocation("Europe/Helsinki") // GMT+2

    return &MemoryManagerServ{
        bringer:      bringer,
        buffer:       NewCircularBuffer(500),
        groupTZ:      groupTZ,
        serverTZ:     serverTZ,
        currentDate:  time.Now().In(groupTZ).Truncate(24 * time.Hour),
        dailyMessages: make([]mlog.Mlog, 0),
    }
}

func (s *MemoryManagerServ) ProcessMessage(msg mlog.Mlog) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Adjust timestamp to group timezone
    msgTime := msg.Timestamp.In(s.groupTZ)
    msgDate := msgTime.Truncate(24 * time.Hour)

    // Check if date has changed
    if !msgDate.Equal(s.currentDate) {
        // Save previous day's messages
        if len(s.dailyMessages) > 0 {
            if err := s.bringer.SaveDailyLog(s.currentDate, s.dailyMessages); err != nil {
                return fmt.Errorf("failed to save daily log: %w", err)
            }
        }
        // Reset for new day
        s.currentDate = msgDate
        s.dailyMessages = make([]mlog.Mlog, 0)
    }

		// Create a copy of msg with adjusted timestamp
    adjustedMsg := msg
    adjustedMsg.Timestamp = msgTime // Устанавливаем время в локальном поясе

    // Add to daily messages
    s.dailyMessages = append(s.dailyMessages, adjustedMsg)

    // Add to circular buffer
    s.buffer.Add(adjustedMsg)

    return nil
}

func (s *MemoryManagerServ) SaveTheDay() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if len(s.dailyMessages) > 0 {
        if err := s.bringer.SaveDailyLog(s.currentDate, s.dailyMessages); err != nil {
            return fmt.Errorf("failed to save daily log: %w", err)
        }
        // Reset daily messages after saving
        s.dailyMessages = make([]mlog.Mlog, 0)
    }

    return nil
}

func (s *MemoryManagerServ) GenerateSummary(count int) (string, error) {
    if count != 100 && count != 200 && count != 500 {
        return "", fmt.Errorf("invalid summary count: %d (must be 100, 200, or 500)", count)
    }

    messages := s.buffer.GetLastN(count)
    if len(messages) == 0 {
        return "", fmt.Errorf("no messages available for summary")
    }

    filePath, err := s.bringer.SaveSummary(messages, count)
    if err != nil {
        return "", fmt.Errorf("failed to save summary: %w", err)
    }

    return filePath, nil
}

