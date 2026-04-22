package repository

type db interface{
	ExecuteSQl(aiText string) error 
	Close() error 
}