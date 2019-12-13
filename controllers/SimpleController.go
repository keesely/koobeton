// SimpleController.go kee > 2019/12/10

package controllers

type SimpleController struct{}

func (c *SimpleController) Get() (int, string) {
	return 200, "simple example"
}
