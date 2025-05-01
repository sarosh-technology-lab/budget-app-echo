package seeders

func RunAllSeeders() {
	SeedRoles()
	SeedPermissions()
	SeedRolePermissions()
	SeedCategories()
	SeedUsers()
}