package settings

const (
	JWTSecret    = "OurSubjectiveJudgmentsWereBiased"
	PasswordSalt = "19thTitleIsLoading#"

	// base url
	BaseURL = "/v1/umg/"

	// database settings
	DBHost       = "DB_HOST"
	DBPort       = "DB_PORT"
	PostgresDB   = "POSTGRES_DB"
	PostgresUser = "POSTGRES_USER"
	PostgresPass = "POSTGRES_PASSWORD"
	DriverName   = "postgres"

	// datetime layouts
	DTLayout     = "2006-01-02T15:04:05"
	UserDTLayout = "Jan 02, 2006 15:04:03"

	// JWT Expiry in minutes
	// Todo: decrease these values
	JWTExpiry            = 2 * 24 * 60
	JWTRefreshExpiry     = 7 * 24 * 60
	RefreshTokenAudience = "https://api.edgecomenergy.ca/v1/umg/token"
)
