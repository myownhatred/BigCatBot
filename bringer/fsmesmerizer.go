package bringer

import (
  "time"
  "os"
  "fmt"
  "path/filepath"
  "encoding/json"

  "Guenhwyvar/lib/mlog"
)

type FileMesmerizer struct {
  logDir string
  summaryDir string
}

func NewFileMesmerizer(logDir, summaryDir string) *FileMesmerizer { 
  return &FileMesmerizer{
    logDir: logDir,
    summaryDir: summaryDir,
  }
}

func (f *FileMesmerizer) SaveDailyLog(date time.Time, messages []mlog.Mlog) error {
  // check for log dir
  if err := os.MkdirAll(f.logDir, 0755); err != nil {
    return fmt.Errorf("failed to create log dir: %w", err)
  }

  //filename template
  filename := filepath.Join(f.logDir, date.Format("2006-01-02")+".json")
  var existingMessages []mlog.Mlog

  // read current file if any and append some logs to it if any
  if _, err := os.Stat(filename); err == nil {
    data, err := os.ReadFile(filename)
    if err != nil {
      return fmt.Errorf("failed to read existing log: %w", err)
    }
    if err := json.Unmarshal(data, &existingMessages); err != nil {
      return fmt.Errorf("failed to unmarshal existing log: %w", err)
    }
  }

  //append new messages
  existingMessages = append(existingMessages, messages...)

  //write back to file
  data, err := json.MarshalIndent(existingMessages, ""," ")
  if err != nil {
    return fmt.Errorf("failed to marshal messages: %w", err)
  }

  if err := os.WriteFile(filename, data, 0644); err != nil {
    return fmt.Errorf("failed to write daily log: %w", err)
  }
  
  return nil
}

func (f *FileMesmerizer) SaveSummary (messages []mlog.Mlog, count int) (string, error) {
  //check for summary dir 
  if err := os.MkdirAll(f.summaryDir, 0755); err != nil {
    return "", fmt.Errorf("there is no summary dir: %w", err)
  }

  //filename for summary 
  filename := filepath.Join(f.summaryDir, fmt.Sprintf("sumy_%d_%d.json", count, time.Now().UnixNano()))

  //write down the stuff
  data, err := json.MarshalIndent(mlog.Msum{Messages:messages}, "", " ")
  if err != nil {
    return "", fmt.Errorf("failed to marshal summary: %w", err)
  }

  if err := os.WriteFile(filename, data, 0644); err != nil {
    return "", fmt.Errorf("failed to write summary to file: %w", err)
  }

  return filename, nil
}



