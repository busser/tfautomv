package pretty

func Error(err error) string {
	return Color(BoxSection("Error", err.Error(), "red"))
	// return Colorf("[red][bold]Error:[reset] %s", err.Error())
}
