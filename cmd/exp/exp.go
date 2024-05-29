package main

import (
	stdctx "context"
	"fmt"

	"github.com/dre4success/lenslocked/context"
	"github.com/dre4success/lenslocked/models"
)


func main() {
	ctx := stdctx.Background()
	
	user := models.User{
		Email: "dre@dre.com",
	}

	ctx = context.WithUser(ctx, &user)
	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)
}
