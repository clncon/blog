package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	"time"
)

type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	Creator   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tag struct {
	BaseModel
	Name    string //标签名称
	IsUsing bool   //是否使用
}
type TagPage struct {
	BaseModel
	TagId  uint
	PageId uint
}
type Page struct {
	BaseModel
	Title       string
	Desc        string
	Body        string `gorm:"type:text"`
	Source      string `gorm:"type.text"`
	IsPublished bool
}

// table users
type User struct {
	gorm.Model
	Email         string    `gorm:"unique_index;default:null"` //邮箱
	Telephone     string    `gorm:"unique_index;default:null"` //手机号码
	Password      string    `gorm:"default:null"`              //密码
	VerifyState   string    `gorm:"default:'0'"`               //邮箱验证状态
	SecretKey     string    `gorm:"default:null"`              //密钥
	OutTime       time.Time `gorm:"default:null"`
	GithubLoginId string    `gorm:"unique_index;default:null"` // github唯一标识
	GithubUrl     string    //github地址
	IsAdmin       bool      //是否是管理员
	AvatarUrl     string    // 头像链接
	NickName      string    // 昵称
	LockState     bool      `gorm:"default:'0'"` //锁定状态
}
type TagCount struct {
	Name  string
	Count int
}

func (tagPage *TagPage) Insert() error {
	return DB.FirstOrCreate(tagPage, "page_id=? and tag_id=?", tagPage.PageId, tagPage.TagId).Error
}
func (tag *Tag) Insert() error {
	return DB.FirstOrCreate(tag, "name = ?", tag.Name).Error
}
func (page *Page) Insert() (uint, error) {
	err := DB.Create(page).Error
	if err != nil {
		return 0, err
	} else {
		var id []uint
		DB.Raw("select last_insert_rowid() as id").Pluck("id", &id)
		return id[0], nil
	}

}
func RemoveTagPageByPageId(pageId string) error {
	return DB.Delete(&TagPage{}, "page_id=?", pageId).Error
}
func GetTagPage(pageId interface{}) ([]*TagPage, error) {
	var tagPages []*TagPage
	rows, err := DB.Raw("select * from tag_pages where page_id=?", pageId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tagPage TagPage
		DB.ScanRows(rows, &tagPage)
		tagPages = append(tagPages, &tagPage)
	}
	return tagPages, nil
}
func GetPage(id interface{}) (*Page, error) {
	var page Page
	err := DB.First(&page, id).Error
	return &page, err
}
func ListTagCount() ([]*TagCount, error) {
	var tagCounts []*TagCount
	rows, err := DB.Raw("select t.name as name,count(*) as count from tags t left join tag_pages tp on t.id=tp.tag_id group by t.name").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tagCount TagCount
		DB.ScanRows(rows, &tagCount)
		logrus.Info(tagCount.Name + "ff")
		tagCounts = append(tagCounts, &tagCount)
	}
	return tagCounts, nil
}
func ListTag() ([]*Tag, error) {
	var tags []*Tag
	rows, err := DB.Raw("select * from tags").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		DB.ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}

/**
 *显示所有的Page
 */
func ListPageAll() ([]*Page, error) {
	var pages []*Page
	rows, err := DB.Raw("select * from pages").Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var page Page
		DB.ScanRows(rows, &page)
		pages = append(pages, &page)
	}
	return pages, nil
}

/**
  显示发布的Page
*/
func ListPage(current, pageSize int) ([]*Page, error) {
	var pages []*Page
	var currentRow = (current - 1) * pageSize
	rows, err := DB.Raw("select * from pages where is_published=?", true).Limit(pageSize).Offset(currentRow).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var page Page
		DB.ScanRows(rows, &page)
		pages = append(pages, &page)
	}
	return pages, nil
}
func Total() (total int) {

	err := DB.Raw("select count(1) from pages where is_published=?", true).Count(&total)
	if err.Error != nil {
		logrus.Error(err.Error)
		total = 0
	}
	return total
}
func DeletePage(id string) error {
	return DB.Delete(Page{}, "id=?", id).Error
}
func UpdatePage(id string, page Page) error {
	return DB.Model(&page).Where("id = ?", id).Updates(
		&Page{Title: page.Title, Desc: page.Desc, Body: page.Body, Source: page.Source, IsPublished: page.IsPublished}).Error
}
func ListTagForIsUsing() ([]*Tag, error) {
	var tags []*Tag
	rows, err := DB.Raw("select * from tags where is_using=?", 1).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		DB.ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}
func MustListTag() []*Tag {
	tags, _ := ListTag()
	return tags
}

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "blog.db")
	//db, err := gorm.Open("mysql", "root:mysql@/wblog?charset=utf8&parseTime=True&loc=Asia/Shanghai")
	if err == nil {
		DB = db
		//db.LogMode(true)
		db.AutoMigrate(&Tag{}, &Page{}, &TagPage{}, &User{})
		db.Model(&TagPage{}).AddUniqueIndex("uk_post_tag", "page_id", "tag_id")
		return db, err
	}
	return nil, err

}

func ListUser() ([]*User, error) {
	var pages []*User
	rows, err := DB.Raw("select * from users").Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		DB.ScanRows(rows, &user)
		pages = append(pages, &user)
	}
	return pages, nil
}
func (user *User) FirstOrCreate() (*User, error) {

	err := DB.FirstOrCreate(user, "github_login_id = ?", user.GithubLoginId).Error
	return user, err
}

// user
// insert user
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.First(&user, "email=?", username).Error
	return &user, err
}
func (user *User) Insert() error {
	return DB.Create(user).Error
}
func GetUser(id interface{}) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	return &user, err
}

func IsGithubIdExists(githubId string, id uint) (*User, error) {
	var user User
	err := DB.First(&user, "github_login_id = ?", user.GithubLoginId).Error
	return &user, err
}
func (user *User) UpdateGithubUserInfo() error {
	var githubLoginId interface{}
	if len(user.GithubLoginId) == 0 {
		githubLoginId = gorm.Expr("NULL")
	} else {
		githubLoginId = user.GithubLoginId
	}

	return DB.Model(user).Update(map[string]interface{}{
		"github_login_id": githubLoginId,
		"avatar_url":      user.AvatarUrl,
		"github_url":      user.GithubUrl,
		"NickName":        user.NickName,
	}).Error

}
func IsExists(email string) (*User, error) {
	var user User
	err := DB.First(&user, "email = ?", email).Error
	return &user, err
}
