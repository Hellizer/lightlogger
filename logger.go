package logger

import (
	"os"
	"sync"
	"time"
)

//Хендл обработчика событий
type LogHandle int

//Освобождение хендла с удалением обработчика
func (lh *LogHandle) Free() {
	h := int(*lh)
	if h >= 0 {
		log.handlres[h] = nil
		*lh = LogHandle(-1)
	}
}

//Делегат для событий лога
type LogEventHandler func(msg *LogMsg)

//TradeLogger основная лог структура
type logger struct {
	handlres    []LogEventHandler // обработчики событий
	level       uint8             // уровень логирования
	ServiceName string            // название сервиса в котором логгер запущен
	msgQue      *linkedQueue      // Очередь сообщений
	waiter      sync.WaitGroup    // ожидание когда все горутины отработают
	locker      sync.Mutex        //
}

//NewTradeLogger создает новый экземляр логгера
func newLogger() *logger {
	l := &logger{}
	l.handlres = make([]LogEventHandler, 1)
	l.handlres[0] = l.defaultHandler
	l.level = 255
	l.msgQue = new(linkedQueue)
	return l
}

//Обработчик по умолчанию - обычный вывод на экран в Stdout и Stderr
func (l *logger) defaultHandler(msg *LogMsg) {
	str := msg.String()
	if msg.Type == LogError {
		os.Stderr.WriteString(str)
		return
	}
	os.Stdout.WriteString(str)
}

func (l *logger) onLogging() {
	l.waiter.Add(1)
	l.locker.Lock()
	msg := l.msgQue.get()
	l.locker.Unlock()
	if msg != nil {
		for i := 0; i < len(l.handlres); i++ {
			l.handlres[i](&msg.LogMsg)
		}
	}

	l.waiter.Done()
}

//Dispose .
func (l *logger) dispose() {
	l.waiter.Wait()
	l.handlres = nil
}

var log = newLogger() //*Logger

//Set Logging level (default = 255)
func SetLogLevel(lvl uint8) {
	log.level = lvl
}

//set name of service
func SetServiceName(str string) {
	log.ServiceName = str
}

//Wait for all out, then free
func SoftDispose() {
	for log.msgQue.cap > 0 {
		time.Sleep(time.Second)
	}

	//log.dispose()
}

//Log добавление записи в лог
func Print(lvl uint8, lmt LogMsgType, obj string, msg string) {
	if lvl <= log.level {
		lmsg := LogMsg{}
		lmsg.TimeAt = time.Now()
		lmsg.ServiceName = log.ServiceName
		lmsg.Level = lvl
		lmsg.Message = msg
		lmsg.ServiceObjectName = obj
		lmsg.Type = lmt
		log.locker.Lock()
		log.msgQue.put(&LogObject{LogMsg: lmsg})
		log.locker.Unlock()
		go log.onLogging()
	}
}

func AddHandler(eh LogEventHandler) LogHandle {
	log.handlres = append(log.handlres, eh)
	return LogHandle(len(log.handlres) - 1)
}
