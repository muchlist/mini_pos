package roles

const (
	RoleOwner    = "owner"
	RoleEmployee = "employee"
	RoleCustomer = "customer"
)

func GetRolesAvailable() []string {
	return []string{RoleOwner, RoleEmployee, RoleCustomer}
}
