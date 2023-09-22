package logger

import (
	"encoding/json"
	"time"
)

//LogMsgType тип сделки
type LogMsgType uint8

const (
	LogNone  LogMsgType = iota // Не определено
	LogError                   // Лог с ошибкой
	LogInfo                    // Лог с инфо
	LogWarn                    // лог с предупреждением
)

//LogMsg сообщение лога
type LogMsg struct {
	Level             uint8      `json:"level"`               // текущий уровнеь записи
	ServiceName       string     `json:"service_name"`        // название сервиса
	ServiceObjectName string     `json:"service_object_name"` // название объекта
	Message           string     `json:"message"`             // логируемое сообщение
	Type              LogMsgType `json:"type"`                // тип сообщения (ошибка, инфо)
	TimeAt            time.Time  `json:"time_stamp"`          // время создания записи
}

type bbuff []byte

func (b *bbuff) writeStrings(str ...string) {
	for _, v := range str {
		*b = append(*b, v...)
	}
}

func (m *LogMsg) String() string {
	var buffer bbuff
	//fmt.Sprintf()
	buffer.writeStrings(m.TimeAt.Format("2006-01-02 15:04:05.000 "))
	switch m.Type {
	case LogNone:
		buffer.writeStrings("[", "NONE", "] ")
	case LogInfo:
		buffer.writeStrings("[", "INFO", "] ")
	case LogWarn:
		buffer.writeStrings("[", "WARN", "] ")
	case LogError:
		buffer.writeStrings("[", "ERRO", "] ")
	}
	buffer.writeStrings("[", m.ServiceName, "] ")
	buffer.writeStrings("[", m.ServiceObjectName, "] ")
	buffer.writeStrings(m.Message, "\n")
	return string(buffer)
}

func (m *LogMsg) ToJsonStr() string {
	jstr, _ := json.Marshal(m)
	return string(jstr)
}

func (m *LogMsg) ToJson() []byte {
	jbyte, _ := json.Marshal(m)
	return jbyte
}
