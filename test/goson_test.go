package goson_test

import (
	"github.com/harlock/goson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

type User struct {
	Name string
	Id int
	Supervisor *User
}

func (this User) Over() bool {
	return this.Id > 100
}

func (this User) Upper() string {
	return strings.ToUpper(this.Name)
}

var _ = Describe("Goson", func() {

	Describe("Creating a JSON object", func() {
		Context("Hash", func() {
			It("Should generage a json object", func() {
				foo := User{Name: "Foo", Id: 11 }
				g := goson.Hash(foo, "Id", "Over()")
				g.Method("Name")
				g.Alias("Uppercase()", "Upper()")
				data, err := g.ToJson()
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`{"id":11,"name":"Foo","over":false,"uppercase":"FOO"}`))
			})
		})

		Context("Array", func() {
			It("Should generage a json array", func() {
				foo := User{Name: "Foo", Id: 1}
				bar := User{Name: "Bar", Id: 200}
				array := []User{ foo, bar }

				g := goson.Array(array, "Id", "Name", "Over()", "Upper()")

				data, err := g.ToJson()
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(`[{"id":1,"name":"Foo","over":false,"upper":"FOO"},{"id":200,"name":"Bar","over":true,"upper":"BAR"}]`))
			})
		})
	})

	Describe("Nested Hash", func() {
		It("Should include the nested hash", func() {
			foo := User{Name: "Foo", Id: 1}
			bar := User{Name: "Bar", Id: 200}
			wadus := User{Name: "Wadus", Id: 300}
			foo.Supervisor = &bar
			bar.Supervisor = &wadus

			g := goson.Hash(foo, "Id")
			supervisor := g.Hash("Supervisor", "Id").Method("Name")
			supervisor.HashAlias("Supervisor", "manager").Alias("manager_id", "Id")

			data, err := g.ToJson()
			Expect(err).To(BeNil())
			Expect(string(data)).To(Equal(`{"id":1,"supervisor":{"id":200,"manager":{"manager_id":300},"name":"Bar"}}`))
		})
	})
})
