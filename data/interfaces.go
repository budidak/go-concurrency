package data

// Normalde main.go içinde application config parametrelerini verirken Models: data.New(db) yazmıştık.
// Bu test aşamasında isteyeceğimiz bir şey değil çünkü database'e bağlanmaya çalışacak.
// Yani testleri çalıştırırken, arkaplanda çalışan bir database isteyecek.
// Çünkü models.go içinde New() fonksiyonuna baktığımızda
// database ile iletişimde olan iki tane struct oluşturup onları Models olarak dönecek.
// Bunu istemiyoruz, database ile doğrudan konuşmasın diye aynı metotları şart koşan interfaces yazdık.
// Artık yeni haliyle, testleri çalıştırdığımızda New() çağrıldığında database ile doğrudan konuşmayan yapılar var elimizde.

// User type satisfies this interface because it has all these methods
type UserInterface interface {
	GetAll() ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetOne(id int) (*User, error)
	Update(user *User) error
	Delete() error
	DeleteByID(id int) error
	Insert(user User) (int, error)
	ResetPassword(password string) error
	PasswordMatches(user *User, plainText string) (bool, error)
}

// Model type satisfies this interface because it has all these methods
type PlanInterface interface {
	GetAll() ([]*Plan, error)
	GetOne(id int) (*Plan, error)
	SubscribeUserToPlan(user User, plan Plan) error
	AmountForDisplay() string
}
