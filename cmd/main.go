package main

import (
	"html/template"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type Contact struct {
	Name  string
	Email string
	Id    int
}

type Data struct {
	Contacts []Contact
}

func NewData() *Data {
	return &Data{
		Contacts: []Contact{
			{
				Name:  "John Doe",
				Email: "john.doe@gmail.com",
				Id:    1,
			},
			{
				Name:  "Jane Doe",
				Email: "jain.doe@gmail.com",
				Id:    2,
			},
		},
	}
}

type FormData struct {
    Errors map[string]string
    Values map[string]string
    Id     int // Add Id field to FormData
}


func NewFormData() FormData {
    return FormData{
        Errors: map[string]string{},
        Values: map[string]string{},
        Id:     0, // Set Id to default value
    }
}


type PageData struct {
	Data      Data
	Form      FormData
	ContactID int // Add ContactID field to PageData
}

func NewContact(id int, name, email string) Contact {
	return Contact{
		Id:    id,
		Name:  name,
		Email: email,
	}
}

func main() {
	e := echo.New()

	data := NewData()
	id := 3

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", PageData{
			Data:      *data,
			Form:      NewFormData(),
			ContactID: 0, // Set ContactID to 0 initially
		})
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if contactExists(data.Contacts, email) {
			formData := FormData{
				Errors: map[string]string{
					"email": "Email already exists",
				},
				Values: map[string]string{
					"name":  name,
					"email": email,
				},
			}

			return c.Render(422, "contact-form", formData)
		}

		contact := NewContact(id, name, email)
		id++
		data.Contacts = append(data.Contacts, contact)

		formData := NewFormData()
		err := c.Render(200, "contact-form", formData)

		if err != nil {
			return err
		}

		return c.Render(200, "oob-contact", contact)
	})

	e.PUT("/contacts/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			return c.String(400, "Id must be an integer")
		}

		name := c.FormValue("name")
		email := c.FormValue("email")

		// Find the contact with the specified id
		var updatedContact *Contact
		for _, contact := range data.Contacts {
			if contact.Id == id {
				contact.Name = name
				contact.Email = email
				updatedContact = &contact
				break
			}
		}

		if updatedContact == nil {
			return c.String(400, "Contact not found")
		}

		// Return the updated contact as JSON
		return c.JSON(200, updatedContact)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			return c.String(400, "Id must be an integer")
		}

		time.Sleep(1 * time.Second)
		deleted := false
		for i, contact := range data.Contacts {
			if contact.Id == id {
				data.Contacts = append(data.Contacts[:i], data.Contacts[i+1:]...)
				deleted = true
				break
			}
		}

		if !deleted {
			return c.String(400, "Contact not found")
		}

		return c.NoContent(200)
	})

	e.Logger.Fatal(e.Start(":42069"))
}

func contactExists(contacts []Contact, email string) bool {
    for _, c := range contacts {
        if c.Email == email {
            return true
        }
    }
    return false // Return false if contact does not exist
}