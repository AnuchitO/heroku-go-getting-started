package main

import (
	"context"
	"errors"
	mlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"gitdev.devops.krungthai.com/aster/ariskill/config"
	"gitdev.devops.krungthai.com/aster/ariskill/database"

	_ "gitdev.devops.krungthai.com/aster/ariskill/docs"
)

// @title						Ariskill API
// @version					1.0
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes					http https
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	cfg := config.C(os.Getenv("ENV"))

	// mlog, graceful := logger.NewZap()
	// defer graceful()

	db, cleanupDBFunc := database.NewMongo(cfg.Database)
	// r := NewRouter(mlog, cfg, db)
	r := NewRouter(nil, nil, db)

	srv := http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	mlog.Println("server start at : " + srv.Addr)

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		cleanupDBFunc()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			mlog.Println("HTTP server Shutdown: " + err.Error())
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		mlog.Fatal("HTTP server ListenAndServe: " + err.Error())
	}

	<-idleConnsClosed
}

func NewRouter(mlog *zap.Logger, cfg *config.Config, db *mongo.Database) *gin.Engine {
	// r := app.NewRouter(mlog)
	r := gin.Default()
	port := os.Getenv("PORT")
	r.GET("/", func(c *gin.Context) {
		msg := "comment config : " + port + " " + os.Getenv("ENV") + " " + os.Getenv("DEV_MONGODB_URI") + " " + time.Now().String()
		c.JSON(http.StatusOK, msg)
	})

	// r.GET("/health", func(c app.Context) {
	// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 	defer cancel()
	// 	if err := db.Client().Ping(ctx, nil); err != nil {
	// 		c.InternalServerError(fmt.Errorf("api server is live: but can't connect to database: %w", err))
	// 		return
	// 	}
	// 	c.OK("ariskill is ready and connected to database")
	// })

	// // packages authen
	// authenHandler := authen.NewAuthenHandler(http.DefaultClient, cfg.GoogleOidc)
	// r.POST("/auth/token", authenHandler.ExchangeForTokens)

	// userStorage := user.NewStorage(db)
	// // packages middlewares
	// r.Use(middlewares.ValidateGoogleIdToken(userStorage.GetOneByEmail, cfg.GoogleOidc, app.RealClock{}))

	// // packages user
	// userHandler := user.NewUserHandler(userStorage)
	// r.GET("/users", userHandler.GetAllUsers)
	// r.GET("/users/email", userHandler.GetUsersData)
	// r.GET("/users/token", userHandler.GetUserByEmailToken)
	// r.GET("/users/hardskills", userHandler.GetHardSkillsMemberByID)

	// // packages profile
	// profileStorage := profile.NewStorage(db)
	// aboutmeUpdateHandler := profile.NewUserHandler(profileStorage)
	// r.PUT("/profile", aboutmeUpdateHandler.UpdateAboutMe)

	// profileHandler := profile.NewProfileHandler(profileStorage)
	// r.GET("/profile", profileHandler.User)

	// skillProfileHandler := profile.NewSkillHandler(profileStorage)
	// r.GET("/profile/skills", skillProfileHandler.GetSkillsByUserID)
	// r.POST("/profile/skills/technical", skillProfileHandler.UpdateTechnicalSkill)
	// r.POST("/profile/skills/soft", skillProfileHandler.UpdateSoftSkill)

	// squadsProfileHandler := profile.NewSquadHandler(profileStorage)
	// r.GET("/profile/squad/:squadID/skill-ratings", squadsProfileHandler.GetUserSkillRatingBySquadID)
	// r.POST("/profile/squad/:squadID/skill-ratings", squadsProfileHandler.RateSkills)

	// // packages membersquad
	// memberSquadStorage := membersquad.NewStorage(db)
	// memberSquadHandler := membersquad.NewMemberSquadHandler(memberSquadStorage)
	// r.PUT("/member-squads/members", memberSquadHandler.AddMemberSquad)
	// r.DELETE("/member-squads/:squadID/members", memberSquadHandler.DeleteMemberSquad)

	// // packages skill
	// skillStorage := skill.NewStorage(db)
	// skillHandler := skill.NewSkillHandler(skillStorage)
	// r.GET("/skills/kind/:kindtype", skillHandler.GetSkillsByKind)
	// r.GET("/skills/:id", skillHandler.SkillByID)
	// r.GET("/hard-skills", skillHandler.SkillByJobRole)

	// // packages squad
	// squadStorage := squad.NewSquadStorage(db)
	// squadMemberStorage := user.NewStorage(db)
	// squadHandler := squad.NewSquadHandler(squadStorage, squadMemberStorage)
	// r.POST("/squads", squadHandler.InsertOneByID)
	// r.GET("/squads", squadHandler.GetAll)
	// r.GET("squads/:squadID", squadHandler.GetOneByID)
	// r.PUT("/squads/:squadID", squadHandler.UpdateOneByID)
	// r.DELETE("/squads/:squadID", squadHandler.DeleteByID)
	// r.GET("/squads/:squadID/member-skills-avg", squadHandler.CalculateSquadMemberAveragePerSkill)
	// r.GET("/squads/:squadID/skills-require-avg", squadHandler.GetAvgSkillRatingByID)

	// // packages cycle
	// cycleStorage := cycle.NewCycleStorage(db)
	// cycleHandler := cycle.NewCycleHandler(cycleStorage)
	// r.POST("/cycles", cycleHandler.InsertOne)
	// r.GET("/cycles/email/user", cycleHandler.GetAllFromUserEmail)
	// r.GET("/cycles/email/:status/:page", cycleHandler.GetAllFromReceiverEmail)
	// r.POST("/cycles/:id", cycleHandler.UpdateByID)
	// r.POST("/cycles/save/:id", cycleHandler.UpdateByIDSave)
	// r.GET("/cycles/:id", cycleHandler.GetOneByID)
	// r.DELETE("/cycles/:id", cycleHandler.DeleteByID)
	// r.GET("/cycles/progress/:id", cycleHandler.GetCycleProgess)
	// r.POST("/cycles/update/:id", cycleHandler.UpdateUserFinalScore)
	// r.PUT("/cycles/goal", cycleHandler.UpdateHardSkillsByEmail)
	// r.GET("/cycles/email/lastest/:email", cycleHandler.GetLatestCycleFromUserEmail)
	// r.GET("/swagger/*any", app.NewSwaggerHandler())
	return r
}
