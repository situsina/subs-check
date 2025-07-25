package interfaces

// Repository 统一的仓库接口
type Repository interface {
	Auth() AuthRepository

	Config() ConfigRepository

	Notify() NotifyRepository
	NotifyTemplate() NotifyTemplateRepository

	Check() CheckRepository

	Sub() SubRepository
	Share() ShareRepository
	SubTemplate() SubTemplateRepository

	Storage() StorageRepository

	Close() error
	Migrate() error
}
