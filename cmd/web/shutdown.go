package main

import (
	"os"
	"os/signal"
	"syscall"
)

// GRACEFUL SHUTDOWN : Main program kapatıldığında, normalde arka planda çalışan goroutine'ler direkt sonlanır. Ama bu goroutine'lerin işlerini tamamlamasını beklemeliyiz ve bu işler tamamladıktan sonra programı sonlandırmalıyız.

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)                      // create new channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // gelen sinyalleri channel'a aktarır => interrupt, terminate
	<-quit                                               // wait for a data from channel; blocks if channel is empty
	app.shutdown()                                       // function call
	os.Exit(0)                                           // exit program with success 0 code
}

func (app *Config) shutdown() {
	app.InfoLog.Println("performing any cleanup tasks...")

	app.Wait.Wait()             // blocks if WaitGroup counter is not zero. (no more mails in the channel)
	app.Mailer.DoneChan <- true // quit the routine
	app.ErrorChanDone <- true

	app.InfoLog.Println("closing channels and shutting down application...")
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)
	close(app.ErrorChan)
	close(app.ErrorChanDone)
}
