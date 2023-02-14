package main

// sendEmail is a helper function that is used for sending a mail to the channel
func (app *Config) sendEmail(msg Message) {
	app.Wait.Add(1)              // increments wait group counter by one.
	app.Mailer.MailerChan <- msg // sends the msg parameter to the mailer channel
}
