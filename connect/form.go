package connect

import (
	"fmt"
	"net/url"
	"strconv"

	"bitbucket.org/goreorto/sqlhero/config"
	"github.com/rivo/tview"
)

// Form ...
type Form struct {
	*tview.Form

	adapter  *tview.DropDown
	name     *tview.InputField
	host     *tview.InputField
	username *tview.InputField
	password *tview.InputField
	database *tview.InputField
	port     *tview.InputField

	test    *tview.Button
	save    *tview.Button
	connect *tview.Button

	OnTest    func(config.Connection)
	OnSave    func(*config.Connection)
	OnConnect func(config.Connection)

	MinWidth   int
	Height     int
	connection *config.Connection
}

// NewForm ...
func NewForm() *Form {
	cf := &Form{
		MinWidth: 34,
		Height:   17,
	}

	u := url.URL{}
	u.User = url.User("")
	cf.connection = &config.Connection{}

	cf.Form = tview.NewForm().
		AddInputField("Name", "", 20, nil, func(txt string) {
			cf.connection.Name = txt
		}).
		AddInputField("Host", "", 20, nil, func(txt string) {
			cf.connection.Host = txt
		}).
		AddInputField("Username", "", 20, nil, func(txt string) {
			cf.connection.Username = txt
		}).
		AddPasswordField("Password", "", 10, '*', func(txt string) {
			cf.connection.Password = txt
		}).
		AddInputField("Database", "", 20, nil, func(txt string) {
			cf.connection.Database = txt
		}).
		AddInputField("Port", "3306", 20, func(txt string, _ rune) bool {
			_, err := strconv.ParseInt(txt, 10, 32)
			return err == nil
		}, func(txt string) {
			port, _ := strconv.ParseInt(txt, 10, 32)
			cf.connection.Port = int(port)
		}).
		AddButton("Test", cf.onTest).
		AddButton("Save", cf.onSave).
		AddButton("Connect", cf.onConnect)

	// caching form items
	cf.name = cf.Form.GetFormItemByLabel("Name").(*tview.InputField)
	cf.host = cf.Form.GetFormItemByLabel("Host").(*tview.InputField)
	cf.username = cf.Form.GetFormItemByLabel("Username").(*tview.InputField)
	cf.password = cf.Form.GetFormItemByLabel("Password").(*tview.InputField)
	cf.database = cf.Form.GetFormItemByLabel("Database").(*tview.InputField)
	cf.port = cf.Form.GetFormItemByLabel("Port").(*tview.InputField)
	cf.test = cf.Form.GetButton(0)
	cf.save = cf.Form.GetButton(1)
	cf.connect = cf.Form.GetButton(2)

	cf.Form.SetBorder(true)
	cf.Form.SetTitle("Connection(C-n)")

	return cf
}

// Clear ...
func (f *Form) Clear() {
	f.name.SetText("")
	f.host.SetText("")
	f.username.SetText("")
	f.password.SetText("")
	f.database.SetText("")
	f.port.SetText("3306")
}

// SetConnection ...
func (f *Form) SetConnection(conn *config.Connection) {
	if conn == nil {
		f.connection = &config.Connection{}
		f.name.SetText("")
		f.host.SetText("")
		f.username.SetText("")
		f.password.SetText("")
		f.database.SetText("")
		f.port.SetText("")
		f.save.SetLabel("Save")

		return
	}

	f.connection = conn
	f.name.SetText(conn.Name)
	f.host.SetText(conn.Host)
	f.username.SetText(conn.Username)
	f.password.SetText(conn.Password)
	f.database.SetText(conn.Database)
	f.port.SetText(fmt.Sprintf("%d", conn.Port))
	f.save.SetLabel("Update")
}

func (f *Form) onTest() {
	if f.OnTest == nil {
		return
	}

	f.OnTest(*f.connection)
}

func (f *Form) onSave() {
	if f.OnSave != nil {
		f.OnSave(f.connection)
	}

	f.save.SetLabel("Update")
}

func (f *Form) onConnect() {
	if f.OnConnect == nil {
		return
	}

	f.OnConnect(*f.connection)
}

// SetRect ...
func (f *Form) SetRect(x, y, width, height int) {
	f.Box.SetRect(x, y, width, height)
	f.Form.SetRect(x+width/2-f.MinWidth/2, y+height/2-f.Height/2, f.MinWidth, f.Height)
}
