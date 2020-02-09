package models

type TestObject interface {
	SetCreatedTime(int64)
	SetLoadedTime(int64)
	SetUpdatedTime(int64)
	SetDeletedTime(int64)
	ClearMetaTimestamps()
}

type TestEntity struct {
	CreatedAt int64
	LoadedAt  int64
	UpdatedAt int64
	DeletedAt int64
}

func (n *TestEntity) SetCreatedTime(time int64) {
	n.CreatedAt = time
}
func (n *TestEntity) SetLoadedTime(time int64) {
	n.LoadedAt = time
}
func (n *TestEntity) SetUpdatedTime(time int64) {
	n.UpdatedAt = time
}
func (n *TestEntity) SetDeletedTime(time int64) {
	n.DeletedAt = time
}
func (n *TestEntity) ClearMetaTimestamps() {
	n.CreatedAt = 0
	n.LoadedAt = 0
	n.UpdatedAt = 0
	n.DeletedAt = 0
}
