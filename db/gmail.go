package legatoDb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/smtp"

	"gorm.io/gorm"
)

const gmailType string = "gmails"

type Gmail struct {
	gorm.Model
	ConnectionID uint
	Connection   *Connection `gorm:"foreignkey:id;references:ConnectionID"`
	Token        string
	Service      Service `gorm:"polymorphic:Owner;"`
}
type GmailLoginData struct {
	To        []string `json:"to"`
	Subject   string   `json:"subject"`
	Password  string   `json:"password"`
	EmailFrom string   `json:"email"`
	Body      string   `json:"body"`
}
type updateGmailData struct {
	ConnectionId uint `json:"connectionId"`
}

func (g *Gmail) String() string {
	return fmt.Sprintf("(@Gmail: %+v)", *g)
}

type SendGmailData struct {
	To     string `json:"to"`
	Messge string `json:"message"`
}

func LoginWithSMTP(emailFrom string, emailPassword string) (smtp.Auth, error) {
	emailHost := "smtp.gmail.com"
	var emailAuth smtp.Auth
	emailAuth = smtp.PlainAuth("", emailFrom, emailPassword, emailHost)
	return emailAuth, nil
}

func SendEmailSmtp(to []string, message string, emailHost string, emailFrom string, subjectmsg string, emailAuth smtp.Auth) (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + subjectmsg + "!\n"
	msg := []byte(subject + mime + "\n" + message)
	addr := fmt.Sprintf("%s:%s", emailHost, "587")
	if err := smtp.SendMail(addr, emailAuth, emailFrom, to, msg); err != nil {
		return false, err
	}
	return true, nil
}
func (ldb *LegatoDB) CreateGmailForScenario(s *Scenario, g Gmail) (*Gmail, error) {
	g.Service.UserID = s.UserID
	g.Service.ScenarioID = &s.ID

	ldb.db.Create(&g)
	ldb.db.Save(&g)

	return &g, nil
}

func (ldb *LegatoDB) UpdateGmail(s *Scenario, servId uint, gn Gmail) error {
	var serv Service
	err := ldb.db.Where(&Service{ScenarioID: &s.ID}).Where("id = ?", servId).Find(&serv).Error
	if err != nil {
		return err
	}

	var g Gmail
	err = ldb.db.Where("id = ?", serv.OwnerID).Preload("Service").Find(&g).Error
	if err != nil {
		return err
	}
	if g.Service.ID != servId {
		return errors.New("the ssh service is not in this scenario")
	}
	var a updateGmailData
	err = json.Unmarshal([]byte(gn.Service.Data), &a)
	if err != nil {
		log.Println("con not update ssh")
	}
	if a.ConnectionId != 0 {
		gn.ConnectionID = a.ConnectionId
		user, _ := ldb.GetUserById(g.Service.UserID)
		con, _ := ldb.GetUserConnectionById(&user, gn.ConnectionID)
		gn.Token = con.Data
	}

	ldb.db.Model(&serv).Updates(gn.Service)
	ldb.db.Model(&g).Updates(gn)

	if gn.Service.ParentID == nil {
		legatoDb.db.Model(&serv).Select("parent_id").Update("parent_id", nil)
	}

	return nil
}

func (ldb *LegatoDB) GetGmailByID(id uint, u *User) (Gmail, error) {
	var g Gmail
	err := ldb.db.Where(&Connection{UserID: u.ID}).Where("ID = ?", id).Find(&g).Error
	if err != nil {
		return Gmail{}, err
	}
	return g, nil
}

func (ldb *LegatoDB) GetGmailByService(serv Service) (*Gmail, error) {
	var g Gmail
	err := ldb.db.Where("id = ?", serv.OwnerID).Preload("Service").Find(&g).Error
	if err != nil {
		return nil, err
	}
	if g.ID != uint(serv.OwnerID) {
		return nil, errors.New("the Gmail service is not in this scenario")
	}

	return &g, nil
}

// Service Interface for Gmail
func (g Gmail) Execute(...interface{}) {
	log.Println("*******Starting Gmail Service*******")
	err := legatoDb.db.Preload("Service").Find(&g).Error
	if err != nil {
		panic(err)
	}
	switch g.Service.SubType {
	case "sendEmail":
		var data GmailLoginData
		err = json.Unmarshal([]byte(g.Service.Data), &data)
		if err != nil {
			log.Print(err)
		}
		emailAuth, _ := LoginWithSMTP(data.EmailFrom, data.Password)
		SendEmailSmtp(data.To, data.Body, "smtp.gmail.com", data.EmailFrom, data.Subject, emailAuth)
	}
	log.Printf("Executing type (%s) : %s\n", gmailType, g.Service.Name)
	g.Next()
}

func (g Gmail) Post() {
	log.Printf("Executing type (%s) node in background : %s\n", gmailType, g.Service.Name)
}

func (g Gmail) Next(...interface{}) {
	err := legatoDb.db.Preload("Service.Children").Find(&g).Error
	if err != nil {
		panic(err)
	}

	log.Printf("Executing \"%s\" Children \n", g.Service.Name)

	for _, node := range g.Service.Children {
		serv, err := node.Load()
		if err != nil {
			log.Println("error in loading services in Next()")
			return
		}
		serv.Execute()
	}

	log.Printf("*******End of \"%s\"*******", g.Service.Name)
}