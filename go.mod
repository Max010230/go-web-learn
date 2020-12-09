module "go-web-learn"

go 1.15

require (
	mircool v0.0.0
)

replace (
	mircool => ./mircool
)