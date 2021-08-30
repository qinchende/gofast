package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

// Redigo初始化的配置参数
type (
	ConnConfig struct {
		Url      string `cnf:",NA"`
		Database string `cnf:",NA"`
		User     string `cnf:",NA"`
		Pass     string `cnf:",NA"`
	}
	MgoX struct {
		DB  *mongo.Database
		Cli *mongo.Client
		Ctx context.Context
	}
)

func (mgX *MgoX) Close() {
	err := mgX.Cli.Disconnect(mgX.Ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("MongoDB %v closed.", mgX.Cli)
}

func NewMongo(cf *ConnConfig) *MgoX {
	url := cf.Url
	pass := cf.Pass
	user := cf.User
	db := cf.Database

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	crt := options.Credential{Username: user, Password: pass, PasswordSet: true, AuthSource: db}
	cOpt := options.Client().ApplyURI(url).SetAuth(crt)
	client, err := mongo.Connect(ctx, cOpt)

	if err != nil {
		log.Fatalf("MongoDB connect err: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("MongoDB ping err: %v", err)
	}

	log.Printf("MongoDB: %s connect sucess!", url)
	return &MgoX{DB: client.Database(db), Cli: client, Ctx: context.Background()}
}
