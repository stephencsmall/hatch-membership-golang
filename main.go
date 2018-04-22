package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/astaxie/beego/orm"
	"github.com/frozzare/go-beeper"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"github.com/wfplhatch/membership/models"
	_ "github.com/wfplhatch/membership/routers"
)

func init() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "./hatch.db")
	orm.RegisterModel(new(models.Person), new(models.Action), new(models.Authorization))
}

func main() {

	beeper.Beep()

	// Set up the ORM
	o := orm.NewOrm()
	o.Using("default")

	// Drop table and re-create.
	force := false
	// Print log.
	verbose := false
	// Error.
	err := orm.RunSyncdb("default", force, verbose)
	if err != nil {
		fmt.Println(err)
	}

	// Set up the interactive shell
	shell := ishell.New()
	shell.SetPrompt("HATCH>")
	shell.Println("HATCH Membership Shell")
	shell.Println("Type 'help' and press ENTER for assistance.")

	// Add some shell commands
	shell.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "add a new member",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)

			newPerson := new(models.Person)

			c.Print("Name: ")
			newPerson.Name = c.ReadLine()

			c.Print("Town: ")
			newPerson.Town = c.ReadLine()

			c.Print("Date of Birth (MM/DD/YYYY): ")
			newPerson.Birthdate, _ = time.Parse("01/02/2006", c.ReadLine())

			c.Print("Library Card Number: ")
			newPerson.LibraryCardNumber, _ = strconv.ParseInt(c.ReadLine(), 10, 64)

			newid, _ := o.Insert(newPerson)
			fmt.Println("Created new user with id ", newid)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "find",
		Help: "Find a member by name",
		Func: func(c *ishell.Context) {
			var people []models.Person
			o.QueryTable("person").Filter("name__icontains", strings.Join(c.Args, " ")).All(&people, "Id", "Name", "Town", "Created", "LibraryCardNumber")
			if len(people) == 0 {
				fmt.Println("No members found.")
			} else {

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{
					"index",
					"HATCH ID Number",
					"Name",
					"Town",
					"Member Since",
					"Library Card Number"})
				for i, v := range people {
					table.Append([]string{
						strconv.Itoa(i),
						strconv.Itoa(v.Id),
						v.Name,
						v.Town,
						v.Created.Format("01/02/2006"),
						strconv.FormatInt(v.LibraryCardNumber, 10)})
				}
				table.Render()
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "link",
		Help: "Link a member to a library card",
		Func: func(c *ishell.Context) {
			c.Print("HATCH ID: ")
			hatchid, err := strconv.Atoi(c.ReadLine())

			c.Print("Library Card Number: ")
			libraryCardNumber, err := strconv.ParseInt(c.ReadLine(), 10, 64)

			if err != nil {
				fmt.Println(err)
			}
			person := models.Person{Id: hatchid}
			person.LibraryCardNumber = libraryCardNumber
			o.Update(&person, "LibraryCardNumber")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "checkin",
		Help: "Check a user in with their library card number",
		Func: func(c *ishell.Context) {
			input, _ := strconv.ParseInt(c.Args[0], 10, 64)
			person := models.Person{LibraryCardNumber: input}
			err := o.Read(&person, "LibraryCardNumber")
			if err == orm.ErrNoRows {
				fmt.Println("Found nothing")
			}

			// Check the age
			age := 0
			o.Raw("SELECT CAST((julianday('now') - julianday(birthdate))/365 AS INTEGER) AS age FROM person WHERE id = ?;", person.Id).QueryRow(&age)

			if age < 18 {
				beeper.Beep(3)
				fmt.Println("This person is a minor! Age: ", age)
			}

			if person.Volunteer == true {
				fmt.Println("This person is a volunteer.")
			}

			//auths := models.Authorization{Person: &person}
			//err := o.Read(&auths,"Target")
			//if err == orm.ErrNoRows {
			//	fmt.Println("This person has no authorizations.")
			//} else {
			//	fmt.Println("This person has the following authorizations:")
			//	for _,v := range auths {
			//		fmt.Println()
			//	}
			//}

			checkin := new(models.Action)
			checkin.Person = &person
			checkin.Entered = true
			o.Insert(checkin)
			fmt.Println("Checking in", person.Name, time.Now().String())
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "authorize",
		Help: "Grant a user an authorization.",
		Func: func(c *ishell.Context) {
			authorizations := []string{
				"3D Printer",
				"Laser Cutter",
				"Sewing Machine",
				"Hand Tools",
				"Soldering"}
			c.Print("HATCH ID: ")
			input, _ := strconv.Atoi(c.ReadLine())
			person := models.Person{Id: input}
			err := o.Read(&person)
			if err == orm.ErrNoRows {
				fmt.Println("Found nothing")
			}

			choice := c.MultiChoice(authorizations, "Choose an authorization:")
			assignedAuth := new(models.Authorization)
			assignedAuth.Person = &person
			assignedAuth.Target = authorizations[choice]
			o.Insert(assignedAuth)

		},
	})

	// Switch between UI and Shell mode here:
	shell.Run()
	//beego.Run()
	//fmt.Println(user.One(&models.Person{},"id","name"))

}
