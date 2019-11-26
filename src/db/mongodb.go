package db
import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Users struct {
	Name string `bson:"name"`
	Age int `bson:"age"`
	Interest []string `bson:"interest"`
}

type LogMgr struct {
	client *mongo.Client
	collection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InItMongodb()  {

	var(
		ctx context.Context
		opts *options.ClientOptions
		client *mongo.Client
		err error
		collection *mongo.Collection
	)
	// 连接数据库
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(G_config.MongodbConnectTimeout) * time.Millisecond)  // ctx
	opts = options.Client().ApplyURI(G_config.MongodbUri)  // opts
	if client, err = mongo.Connect(ctx,opts); err != nil{
		fmt.Println(err)
		return
	}

	//链接数据库和表
	collection = client.Database("mygood").Collection("mytest")

	//赋值单例
	G_logMgr = &LogMgr{
		client:client,
		collection:collection,
	}

	//G_logMgr.SaveMongodb()
	//G_logMgr.UpdateMongo()
	//G_logMgr.SelectMongodb()
	//G_logMgr.DeleteMongo()
}


//保存数据
func (logMgr *LogMgr)SaveMongodb() (err error) {
	var(
		insetRest *mongo.InsertOneResult
		id interface{}
		users []interface{}
		insertManRest *mongo.InsertManyResult
	)

	//TODO: 单个写入
	//userId := bson.NewObjectId().String()
	user := Users{"liulong", 23, []string{"爬山", "敲代码"}}
	if insetRest, err = logMgr.collection.InsertOne(context.TODO(), &user); err != nil{
		fmt.Println(err)
		return
	}
	id = insetRest.InsertedID
	fmt.Println(id)

	//TODO: 批量写入
	users = append(users, &Users{Name:"zhangsan", Age:25}, &Users{Name:"lisi", Age:33}, &Users{Name:"wangwu", Age:66})
	if insertManRest, err = logMgr.collection.InsertMany(context.TODO(), users); err != nil{
		fmt.Println(err)
		return
	}
	for _, v := range insertManRest.InsertedIDs{
		fmt.Println(v)
	}

	return
}


//查询数据
func (logMgr *LogMgr)SelectMongodb() (err error)  {

	var(
		cur *mongo.Cursor
		user *Users
		ctx context.Context

	)
	ctx = context.TODO()

	//TODO: 单个查询

	//if err = logMgr.collection.FindOne(ctx, bson.D{{"name", "liulong"}}).Decode(&user); err != nil{
	//	fmt.Println(err)
	//	return
	//}
	//logMgr.collection.FindOne(ctx, bson.M{"name":"liulong", "age": 23}).Decode(&user)
	//if cur, err = logMgr.collection.Find(ctx, bson.D{{"name", "liulong"}, {"age", 23}}); err !=nil{
	//	fmt.Println("此处错误", err)
	//	return
	//}

	if err = logMgr.collection.FindOne(ctx, bson.M{"name":"liulong"}).Decode(&user); err != nil{
		fmt.Println(err)
		return
	}

	fmt.Println(user)


	//TODO：多个查询； 查找age>=25的，只显示3个，从大到小排序
	if cur, err = logMgr.collection.Find(ctx, bson.M{"age":bson.M{"$gte":25}}, options.Find().SetLimit(3), options.Find().SetSort(bson.M{"age":-1})); err != nil{
		fmt.Println(err)
		return
	}
	defer cur.Close(ctx)
	//
	for cur.Next(ctx){
		//TODO: 第一种解析方法，解析后是结构体
		user = &Users{}
		if err = cur.Decode(user); err != nil{
			fmt.Println(err)
		}
		fmt.Println(user)

		//	//TODO: 第二种解析方法，解析后是map
		var result bson.M
		if err = cur.Decode(&result); err != nil{
			fmt.Println(err)
		}
		fmt.Println(result["_id"], result)

	}

	return
}


//更新数据
func (logMgr *LogMgr)UpdateMongo() (err error)  {

	var(
		ctx context.Context
		updateRet *mongo.UpdateResult
	)

	//更新name是liulong33的用户，没有就创建
	//logMgr.collection.UpdateOne(ctx, bson.M{"name": "liulong33"}, bson.M{"$set": bson.M{"age": 78}}, options.Update().SetUpsert(true))
	if updateRet, err = logMgr.collection.UpdateOne(ctx, bson.M{"name": "liulong33"}, bson.M{"$set": bson.M{"age": 78}}); err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("更新的个数", updateRet.ModifiedCount)
	return
}


//删除数据
func (logMgr *LogMgr)DeleteMongo()(err error)  {

	var(
		ctx context.Context
		delRet *mongo.DeleteResult
	)
	if delRet, err = logMgr.collection.DeleteMany(ctx, bson.M{"age":bson.M{"$gte":10}}); err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(delRet.DeletedCount)
	return
}