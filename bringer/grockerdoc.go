package bringer

// Пакет grokker реализует интеграцию с Grok API (xAI).
//
// # Архитектура
//
// Три слоя:
//
//  1. GrokClient — транспортный слой.
//     Один клиент = один API-ключ = один плательщик.
//     Единственная точка отправки HTTP: doRequest().
//
//  2. Conversation — история диалога для конкретного чата.
//     Хранит system prompt + историю user/assistant сообщений.
//     Используется командами /grokc и /grokcd.
//
//  3. GrokGrokker — фасад. Единственный тип используемый снаружи (в servitor).
//     Содержит два именованных клиента и хранилище разговоров.
//
// # Модель биллинга
//
// Разные функции оплачиваются разными людьми:
//
//   - Grik (MothershipKey)  — общие запросы, диалоги, анализ чата
//   - GrikDND (DNDKey)      — генерация DND-биографий и параметрические запросы
//
// Если DNDKey не задан в конфиге — NewGrokker создаст клиент с пустым ключом
// и первый же запрос упадёт с 401. Лучше добавить фолбэк или проверку в Validate().
//
// # Схема вызовов
//
// Публичные методы GrokGrokker и их цепочки:
//
//	SimpleAnswer(prompt)
//	  └── Grik.callGenericGrokAPI(prompt)
//	        └── doRequest(быстрая модель, весёлый стиль)
//	              └── formatUsage(resp) → строка с токенами и стоимостью
//
//	DNDBiogen(prompt)
//	  └── GrikDND.callBiogenGrokAPI(prompt)
//	        ├── tags() → 8 случайных тегов для system prompt
//	        └── doRequest(полная модель, стиль писателя-юмориста)
//	              └── formatUsage(resp)
//
//	GenGrok(prompt, role, temp)
//	  └── GrikDND.callParGrokAPI(prompt, role, temp)
//	        └── doRequest(полная модель, роль и температура снаружи)
//	              └── formatUsage(resp)
//
//	SendMessageInConversation(chatID, content)
//	  ├── GetOrCreateConversation(chatID) → Conversation  [mu.Lock]
//	  ├── conv.AppendMessage("user", content)
//	  ├── Grik.doRequest(все сообщения из истории + настройки conv)
//	  └── conv.AppendMessage("assistant", ответ)
//
//	ClearConversation(chatID)   — сбросить историю, сохранить system prompt  [mu.Lock]
//	DeleteConversation(chatID)  — удалить разговор полностью                 [mu.Lock]
//
//	AnalChatDay(filename)
//	  ├── os.ReadFile(filename) → []byte
//	  ├── json.Unmarshal → []mlog.Mlog
//	  └── Grik.AnalyzeChat(messages)
//	        └── Grik.AnalyzeChatSummary(messages)
//	              ├── formatMessagesForAnalysis(messages, 3000) → текст для промпта
//	              └── callGrokAPI(prompt)   ← UwU-стиль анализатора
//	                    └── doRequest(полная модель, температура 0.5)
//
// # Потокобезопасность
//
// GrokGrokker.mu защищает доступ к convs и chatClients.
// GrokClient потокобезопасен — http.Client и APIKey не меняются после создания.
// Conversation НЕ потокобезопасна сама по себе — защищается через mu GrokGrokker.
//
// # Известные проблемы
//
//   - fmt.Printf/Println в AnalyzeChat (логи без slog, видны только в stdout)
//   - AnalysisResult.TopicStarters всегда пустая строка — поле зарезервировано
//   - Conversation не имеет TTL: старые сессии живут до перезапуска бота
//     (TODO: очистка через /grokcd или автоматический TTL)
//   - Тарифы в formatUsage захардкожены на момент написания:
//     input $3.00/1M, cached $0.75/1M, output $15.00/1M
