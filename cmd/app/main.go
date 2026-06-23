package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/bookshop/internal/service/ollama"
	"github.com/bookshop/internal/service/ollama/embedding"
	llm2 "github.com/bookshop/internal/service/ollama/llm"
	"github.com/bookshop/internal/storage/knowledge_base"

	"github.com/bookshop/internal/logger"
	adminS "github.com/bookshop/internal/service/admin"
	acS "github.com/bookshop/internal/service/ai_chat"
	adS "github.com/bookshop/internal/service/ai_dialog"
	bookS "github.com/bookshop/internal/service/book"
	brS "github.com/bookshop/internal/service/book_revs"
	bvS "github.com/bookshop/internal/service/book_views"
	cartS "github.com/bookshop/internal/service/cart"
	ciS "github.com/bookshop/internal/service/cart_items"
	favS "github.com/bookshop/internal/service/fav"
	fiS "github.com/bookshop/internal/service/fav_items"
	aiservice "github.com/bookshop/internal/service/knowledge_base"
	oiS "github.com/bookshop/internal/service/order_items"
	orderS "github.com/bookshop/internal/service/orders"
	pageS "github.com/bookshop/internal/service/page"
	pbS "github.com/bookshop/internal/service/physical_books"
	readS "github.com/bookshop/internal/service/reading"
	rsS "github.com/bookshop/internal/service/reading_sessions"
	recS "github.com/bookshop/internal/service/recommendation"
	statsS "github.com/bookshop/internal/service/stats"
	spayS "github.com/bookshop/internal/service/subscription_payments"
	splanS "github.com/bookshop/internal/service/subscription_plans"
	userS "github.com/bookshop/internal/service/user"
	usubS "github.com/bookshop/internal/service/user_subscriptions"
	ubS "github.com/bookshop/internal/service/users_books"
	"github.com/bookshop/internal/storage"
	"github.com/bookshop/internal/storage/admin"
	"github.com/bookshop/internal/storage/ai_chat"
	"github.com/bookshop/internal/storage/book"
	"github.com/bookshop/internal/storage/book_revs"
	bookviews "github.com/bookshop/internal/storage/book_views"
	"github.com/bookshop/internal/storage/cart"
	"github.com/bookshop/internal/storage/cart_items"
	"github.com/bookshop/internal/storage/fav"
	"github.com/bookshop/internal/storage/fav_items"
	"github.com/bookshop/internal/storage/order_items"
	"github.com/bookshop/internal/storage/orders"
	"github.com/bookshop/internal/storage/page"
	"github.com/bookshop/internal/storage/physical_books"
	"github.com/bookshop/internal/storage/reading"
	"github.com/bookshop/internal/storage/reading_sessions"
	"github.com/bookshop/internal/storage/recommendation"
	redislocal "github.com/bookshop/internal/storage/redis"
	redisHistory "github.com/bookshop/internal/storage/redis/history"
	"github.com/bookshop/internal/storage/stats"
	subscription_payments "github.com/bookshop/internal/storage/subscription_payments"
	subscription_plans "github.com/bookshop/internal/storage/subscription_plans"
	"github.com/bookshop/internal/storage/user"
	user_subscriptions "github.com/bookshop/internal/storage/user_subscriptions"
	"github.com/bookshop/internal/storage/users_books"
	adminH "github.com/bookshop/internal/transport/web/controllers/admin"
	acH "github.com/bookshop/internal/transport/web/controllers/ai_chat"
	aihandler "github.com/bookshop/internal/transport/web/controllers/ai_service"
	bookH "github.com/bookshop/internal/transport/web/controllers/book"
	brH "github.com/bookshop/internal/transport/web/controllers/book_revs"
	bvH "github.com/bookshop/internal/transport/web/controllers/book_views"
	cartH "github.com/bookshop/internal/transport/web/controllers/cart"
	ciH "github.com/bookshop/internal/transport/web/controllers/cart_items"
	favH "github.com/bookshop/internal/transport/web/controllers/fav"
	fiH "github.com/bookshop/internal/transport/web/controllers/fav_items"
	oiH "github.com/bookshop/internal/transport/web/controllers/order_items"
	orderH "github.com/bookshop/internal/transport/web/controllers/orders"
	pageH "github.com/bookshop/internal/transport/web/controllers/page"
	pbH "github.com/bookshop/internal/transport/web/controllers/physical_books"
	readH "github.com/bookshop/internal/transport/web/controllers/reading"
	rsH "github.com/bookshop/internal/transport/web/controllers/reading_sessions"
	recH "github.com/bookshop/internal/transport/web/controllers/recommendation"
	statsH "github.com/bookshop/internal/transport/web/controllers/stats"
	spayH "github.com/bookshop/internal/transport/web/controllers/subscription_payments"
	splanH "github.com/bookshop/internal/transport/web/controllers/subscription_plans"
	userH "github.com/bookshop/internal/transport/web/controllers/user"
	usubH "github.com/bookshop/internal/transport/web/controllers/user_subscriptions"
	ubH "github.com/bookshop/internal/transport/web/controllers/users_books"
	"github.com/bookshop/internal/transport/web/middleware"
	adminR "github.com/bookshop/internal/transport/web/routers/admin"
	acR "github.com/bookshop/internal/transport/web/routers/ai_chat"
	airouter "github.com/bookshop/internal/transport/web/routers/ai_service"
	bookR "github.com/bookshop/internal/transport/web/routers/book"
	brR "github.com/bookshop/internal/transport/web/routers/book_revs"
	bvR "github.com/bookshop/internal/transport/web/routers/book_views"
	cartR "github.com/bookshop/internal/transport/web/routers/cart"
	ciR "github.com/bookshop/internal/transport/web/routers/cart_items"
	favR "github.com/bookshop/internal/transport/web/routers/fav"
	fiR "github.com/bookshop/internal/transport/web/routers/fav_items"
	oiR "github.com/bookshop/internal/transport/web/routers/order_items"
	orderR "github.com/bookshop/internal/transport/web/routers/orders"
	pageR "github.com/bookshop/internal/transport/web/routers/page"
	pbR "github.com/bookshop/internal/transport/web/routers/physical_books"
	readR "github.com/bookshop/internal/transport/web/routers/reading"
	rsR "github.com/bookshop/internal/transport/web/routers/reading_sessions"
	recR "github.com/bookshop/internal/transport/web/routers/recommendation"
	statsR "github.com/bookshop/internal/transport/web/routers/stats"
	spayR "github.com/bookshop/internal/transport/web/routers/subscription_payments"
	splanR "github.com/bookshop/internal/transport/web/routers/subscription_plans"
	userR "github.com/bookshop/internal/transport/web/routers/user"
	usubR "github.com/bookshop/internal/transport/web/routers/user_subscriptions"
	ubR "github.com/bookshop/internal/transport/web/routers/users_books"
	"github.com/bookshop/internal/transport/web/server"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: .env file not found:", err.Error())
	}

	if err := middleware.InitFirebaseAuthWithProject("bookshopauth-38202"); err != nil {
		fmt.Println(err)
		return
	}

	loggerSv := logger.New(slog.LevelDebug)
	loggerSv.Info("Application starting...")

	lstor := loggerSv.With("Layer", "Storage")
	lserv := loggerSv.With("Layer", "Service")
	lhand := loggerSv.With("Layer", "Handlers")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"))

	loggerSv.Info("Connecting to database...", "connection_string", connStr)

	db, err := storage.Connection(connStr)
	if err != nil {
		loggerSv.Error("Failed to connect to database", "error", err)
		panic(err)
	}
	defer db.Close()
	loggerSv.Info("Database connected successfully")

	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisOptions := &redis.Options{Addr: addr, Password: os.Getenv("REDIS_PASSWORD"), DB: 0}
	redisObj, err := redislocal.Connection(redisOptions)
	if err != nil {
		loggerSv.Error("Failed to connect to redis", "error", err)
		panic(err)
	}
	defer redisObj.Close()
	loggerSv.Info("Redis connected successfully")

	hst := redisHistory.New(redisObj, lstor)

	loggerSv.Info("Initializing user components...")
	ust := user.New(db, lstor)
	us := userS.New(ust, lserv)
	uh := userH.New(us, lhand)
	ur := userR.New(uh)

	loggerSv.Info("Initializing book components...")
	bst := book.New(db, lstor)
	bs := bookS.New(bst, ust, lserv)
	bh := bookH.New(bs, lhand)
	br := bookR.New(bh)

	loggerSv.Info("Initializing admin components...")
	admst := admin.New(db, lstor)
	admserv := adminS.New(admst, lserv)
	admh := adminH.New(admserv, lhand)
	admr := adminR.New(admh)

	loggerSv.Info("Initializing physical_books components...")
	pbst := physical_books.New(db, lstor)
	pbserv := pbS.New(pbst, lserv)
	pbh := pbH.New(pbserv, lhand)
	pbr := pbR.New(pbh)

	loggerSv.Info("Initializing page components...")
	pst := page.New(db, lstor)
	ps := pageS.New(pst, lserv)
	ph := pageH.New(ps, lhand)
	pr := pageR.New(ph)

	loggerSv.Info("Initializing users_books components...")
	ubst := users_books.New(db, lstor)
	ubs := ubS.New(ubst, ust, lserv)
	ubh := ubH.New(ubs, lhand)
	ubr := ubR.New(ubh)

	loggerSv.Info("Initializing cart components...")

	cst := cart.New(db, lstor)
	cserv := cartS.New(cst, ust, lserv)
	ch := cartH.New(cserv, lhand)
	cr := cartR.New(ch)

	loggerSv.Info("Initializing cart_items components...")

	cist := cart_items.New(db, lstor)
	ciserv := ciS.New(cist, cst, ust, lserv)
	cih := ciH.New(ciserv, lhand)
	cir := ciR.New(cih)

	loggerSv.Info("Initializing fav components...")

	fst := fav.New(db, lstor)
	fserv := favS.New(fst, ust, lserv)
	fh := favH.New(fserv, lhand)
	fr := favR.New(fh)

	loggerSv.Info("Initializing fav_items components...")

	fist := fav_items.New(db, lstor)
	fiserv := fiS.New(fist, fst, ust, lserv)
	fih := fiH.New(fiserv, lhand)
	fir := fiR.New(fih)

	loggerSv.Info("Initializing reading components...")

	readst := reading.New(db, lstor)
	readserv := readS.New(readst, ust, bst, hst, lserv)
	readh := readH.New(readserv, lhand)
	readr := readR.New(readh)

	loggerSv.Info("Initializing reading_sessions components...")

	rsst := reading_sessions.New(db, lstor)
	rsserv := rsS.New(rsst, bst, hst, ust, lserv)
	rsh := rsH.New(rsserv, lhand)
	rsr := rsR.New(rsh)

	loggerSv.Info("Initializing book_revs components...")

	brst := book_revs.New(db, lstor)
	brserv := brS.New(brst, ust, lserv)
	brh := brH.New(brserv, lhand)
	brr := brR.New(brh)

	loggerSv.Info("Initializing book_views components...")

	bvst := bookviews.New(db, lstor)
	bvserv := bvS.New(bvst, bst, hst, ust, lserv)
	bvh := bvH.New(bvserv, lhand)
	bvr := bvR.New(bvh)

	loggerSv.Info("Initializing stats components...")

	sst := stats.New(db, lstor)
	sserv := statsS.New(sst, ust, lserv)
	sh := statsH.New(sserv, lhand)
	sr := statsR.New(sh)

	loggerSv.Info("Initializing recommendation components...")

	recst := recommendation.New(db, lstor)
	recserv := recS.New(recst, ust, lserv)
	rech := recH.New(recserv, lhand)
	recr := recR.New(rech)

	loggerSv.Info("Initializing user_subscriptions components...")

	usubst := user_subscriptions.New(db, lstor)
	usubserv := usubS.New(usubst, ust, lserv)
	usubh := usubH.New(usubserv, lhand)
	usubr := usubR.New(usubh)

	loggerSv.Info("Initializing subscription_payments components...")

	spayst := subscription_payments.New(db, lstor)
	spayserv := spayS.New(spayst, ust, lserv)
	spayh := spayH.New(spayserv, lhand)
	spayr := spayR.New(spayh)

	loggerSv.Info("Initializing subscription_plans components...")

	splanst := subscription_plans.New(db, lstor)
	splanserv := splanS.New(splanst, lserv)
	splanh := splanH.New(splanserv, lhand)
	splanr := splanR.New(splanh)

	loggerSv.Info("Initializing orders components...")

	ost := orders.New(db, lstor)
	oserv := orderS.New(ost, ust, lserv)
	oh := orderH.New(oserv, lhand)
	or := orderR.New(oh)

	loggerSv.Info("Initializing order_items components...")

	oist := order_items.New(db, lstor)
	oiserv := oiS.New(oist, ust, lserv)
	oih := oiH.New(oiserv, lhand)
	oir := oiR.New(oih)

	loggerSv.Info("Initializing ai_chat components...")

	acst := ai_chat.New(db, lstor)
	acserv := acS.New(acst, ust, lserv)

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	ollamaClient := ollama.New(ollamaURL)
	embedder := embedding.New(ollamaClient)
	llm := llm2.New(ollamaClient)

	aist := knowledge_base.New(db, lstor)
	aiserv := aiservice.New(aist, embedder, llm, lserv)
	adserv := adS.New(acserv, aiserv, lserv)

	ach := acH.New(acserv, adserv, lhand)
	acr := acR.New(ach)

	aihand := aihandler.New(aiserv, lhand)
	airout := airouter.New(aihand)

	loggerSv.Info("Creating server...")
	server := server.New(ur, br, pr, ubr, cir, cr, fir, fr, readr, rsr, brr, bvr, sr, recr, usubr, spayr, splanr, pbr, or, oir, acr, airout, admr)
	router := server.Create()

	loggerSv.Info("Server starting", "port", ":3000")
	loggerSv.Info("Available endpoints:")

	if err := router.Run(":3000"); err != nil {
		loggerSv.Error("Server failed to start", "error", err)
		panic(err)
	}
}
