package main

/*
Attach endpoints to handlerfunctions
*/
func (a application) Routes() {

	app.router.HandleFunc("/create-recipe", app.createRecipeHandler)
	app.router.HandleFunc("/browse", app.allRecipePage)
	app.router.HandleFunc("/myrecipes", app.readRecipe)
	app.router.HandleFunc("/update-recipe/{_id}/", app.updateRecipe)
	app.router.HandleFunc("/remove-recipe{_id}/", app.deleteRecipeHandler)
	app.router.HandleFunc("/support", app.supportHandler)
	app.router.HandleFunc("/mypantry", app.pantryHandler)
	app.router.HandleFunc("/login", app.loginHandler)
	app.router.HandleFunc("/register", app.registerHandler)
	app.router.HandleFunc("/", app.showHandler)
	app.router.HandleFunc("/profile", app.profileHandler)
	app.router.HandleFunc("/logout", app.logoutHandler)
	app.router.HandleFunc("/about", app.aboutHandler)
}
