package main 
import ("fmt" ; "net/http" ; "encoding/json" ; "time" ; "context" ; "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" ; "go.mongodb.org/mongo-driver/bson/primitive" ; "github.com/golang-jwt/jwt" ; "strings" ; "os" )

var client *mongo.Client
var apicollection *mongo.Collection
var checkcollection *mongo.Collection
var usercollection *mongo.Collection




type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username         string             `bson:"username" json:"username"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"password"`
}


type API struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"userId" json:"userId"`

	Name             string             `bson:"name" json:"name"`
	URL              string             `bson:"url" json:"url"`

	LastStatus       string             `bson:"lastStatus" json:"lastStatus"`
	LastStatusCode   int                `bson:"lastStatusCode" json:"lastStatusCode"`
	LastResponseTime int64              `bson:"lastResponseTime" json:"lastResponseTime"`
	LastCheckedAt    time.Time          `bson:"lastCheckedAt" json:"lastCheckedAt"`

	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
}


type Check struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	APIID            primitive.ObjectID `bson:"apiId" json:"apiId"`

	Status           string             `bson:"status" json:"status"`
	StatusCode       int                `bson:"statusCode" json:"statusCode"`
	ResponseTime     int64              `bson:"responseTime" json:"responseTime"`

	CheckedAt        time.Time          `bson:"checkedAt" json:"checkedAt"`
}

type loginuser struct {
	Email   string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type apiRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var jwtSecret = []byte("ma_cle_secrete")

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next(w, r)
    }
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ;
	}
	var user User;
	err := json.NewDecoder(r.Body).Decode(&user);
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest);
		return ;
	}


	count , err := usercollection.CountDocuments(context.Background(), bson.M{"username": user.Username , "email": user.Email , "password": user.Password});
	if err != nil {
		http.Error(w, "Failed to check existing user", http.StatusInternalServerError)
		return ;
	}
	
	if count > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return ;
	}
	result, err := usercollection.InsertOne(context.Background(), user);
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError);
		return ;
	}

	fmt.Println("User created with ID: ", result.InsertedID);

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ;
	}
	var loginUser loginuser;
	err := json.NewDecoder(r.Body).Decode(&loginUser);
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return ;
	}
	fmt.Println("Login attempt for email: ", loginUser.Email);
	fmt.Println("Login attempt for password: ", loginUser.Password);

	var user User;
	err = usercollection.FindOne(context.Background(), bson.M{"email": loginUser.Email}).Decode(&user);

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return ;
	}
	fmt.Println("User password: ", user.Password);
	fmt.Println("Login password: ", loginUser.Password);

	if user.Password != loginUser.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return ;
	}
	


	claims := jwt.MapClaims{
		"userId": user.ID ,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return ;
	}


	loginResp := loginResponse{
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(loginResp);	


}


func getapis(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ;
	}

	authHeader := r.Header.Get("Authorization");
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return ;
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ");
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    	return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return ;
	}

	claims := token.Claims.(jwt.MapClaims);
	userIdHex := claims["userId"].(string)

	userObjectID, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	cursor, err := apicollection.Find(
		context.Background(),
		bson.M{"userId": userObjectID},
	)
	if err != nil {
		http.Error(w, "Failed to fetch APIs", http.StatusInternalServerError)
		return ;
	}
	var apis []API;
	for cursor.Next(context.Background()) {
		var api API;
		err := cursor.Decode(&api);
		if err != nil {
			http.Error(w, "Failed to decode API", http.StatusInternalServerError)
			return ;
		}
		apis = append(apis, api);
	}

	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(apis);
}


func createAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ;
	}
	var apiReq apiRequest;
	err := json.NewDecoder(r.Body).Decode(&apiReq);
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return ;
	}
	fmt.Println("Creating API with Name: ", apiReq.Name);
	fmt.Println("Creating API with URL: ", apiReq.URL);

	authHeader := r.Header.Get("Authorization");
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return ;
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ");
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return ;
	}

	claims := token.Claims.(jwt.MapClaims);
	userId := claims["userId"].(string);
	var api API;
	api.Name = apiReq.Name;
	api.URL = apiReq.URL;
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	api.UserID = objectID

	result, err := apicollection.InsertOne(context.Background(), api);
	if err != nil {
		http.Error(w, "Failed to create API", http.StatusInternalServerError)
		return ;
	}
	fmt.Println("API created with ID: ", result.InsertedID);
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		createAPI(w, r);
	}else if r.Method == "GET" {
		getapis(w, r);
	}
}


func startMonitoring() {
	ticker := time.NewTicker(1 * time.Minute);
	defer ticker.Stop();
	for {
		<- ticker.C
		var apis []API;
		cursor , err := apicollection.Find(context.Background(), bson.M{})
		if err != nil {
			fmt.Println("Error fetching APIs: ", err);
			continue ;
		}
		for cursor.Next(context.Background()) {
			var api API;
			err := cursor.Decode(&api);
			apis = append(apis, api);
			if err != nil {
				fmt.Println("Error decoding API: ", err);
				continue ;
			}
		}

		for _, api := range apis {
			checkAPI(api);
		}
	}
}


func checkAPI(api API) {
	start := time.Now();
	resp , err := http.Get(api.URL);
	responseTime := time.Since(start).Milliseconds();
	check := Check{
		APIID:        api.ID,
		ResponseTime: responseTime,
		CheckedAt:    time.Now(),
	}

	if err != nil {
		check.Status = "DOWN"
		check.StatusCode = 0 ;
	}else{
		defer resp.Body.Close();
		check.StatusCode = resp.StatusCode
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			check.Status = "UP"
		}else{
			check.Status = "DOWN"
		}
	}


	result , err := checkcollection.InsertOne(context.Background(), check);
	if err != nil {
		fmt.Println("Error inserting check result: ", err);
		return ;
	}
	fmt.Println("Check result inserted with ID: ", result.InsertedID);

	update := bson.M{
		"$set": bson.M{
			"lastStatus":       check.Status,
			"lastStatusCode":   check.StatusCode,
			"lastResponseTime": check.ResponseTime,
			"lastCheckedAt":    check.CheckedAt,
		},
	}

	_, err = apicollection.UpdateByID(context.Background(), api.ID, update)
	if err != nil {
		fmt.Println("Error updating API: ", err);
	}
}

func deleteAPI(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ;
	}

	idstr := strings.TrimPrefix(r.URL.Path, "/dashboard/")
	if idstr == "" {
		http.Error(w, "API ID is required", http.StatusBadRequest)
		return ;
	}
	id, err := primitive.ObjectIDFromHex(idstr);
	if err != nil {
		http.Error(w, "Invalid API ID", http.StatusBadRequest)
		return ;
	}

	authHeader := r.Header.Get("Authorization");
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return ;
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ");
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return ;
	}

	claims := token.Claims.(jwt.MapClaims);
	userId := claims["userId"].(string);
	objectID, err := primitive.ObjectIDFromHex(userId);
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return ;
	}	

	 result , err := apicollection.DeleteOne(context.Background(), bson.M{"_id": id, "userId": objectID});
	if err != nil {
		http.Error(w, "Failed to delete API", http.StatusInternalServerError)
		return ;
	}

	if result.DeletedCount == 0 {
		http.Error(w, "API not found or not authorized", http.StatusNotFound)
		return ;
	}

	w.WriteHeader(http.StatusOK);

}


func getChecks(w http.ResponseWriter, r *http.Request) {

	 fmt.Println("getChecks called")
    fmt.Println("Method:", r.Method)
    fmt.Println("Path:", r.URL.Path)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return ;
	}

	idstr := strings.TrimPrefix(r.URL.Path, "/dashboard/api/");
	idstr = strings.TrimSuffix(idstr, "/history");
	id , err := primitive.ObjectIDFromHex(idstr);

	if err != nil {
		http.Error(w, "Invalid API ID", http.StatusBadRequest)
		return ;
	}

	cursor , err := checkcollection.Find(context.Background(), bson.M{"apiId": id});
	if err != nil {
		http.Error(w, "Failed to fetch checks", http.StatusInternalServerError)
		return ;
	}

	var checks []Check;
	for cursor.Next(context.Background()) {
		var check Check ;
		err := cursor.Decode(&check);
		if err != nil {
			http.Error(w, "Failed to decode check", http.StatusInternalServerError)
			return ;
		}
		checks = append(checks, check);
	}

	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(checks);

}

func main() {
	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second);
	defer cancel();
	var err error;
	URI := os.Getenv("MONGODB_URI");
	if URI == "" {
		fmt.Println("MONGODB_URI is not set");
	}
	client , err  = mongo.Connect(ctx , options.Client().ApplyURI(URI));

	if err != nil {
		fmt.Println("Error connecting to MongoDB: " , err);
		return ;
	}

	apicollection = client.Database("amd").Collection("apicollection");
	checkcollection = client.Database("amd").Collection("checkcollection");
	usercollection = client.Database("amd").Collection("usercollection");
	

	http.HandleFunc("/signup", enableCORS(createUser));
	http.HandleFunc("/login", enableCORS(login));
	http.HandleFunc("/dashboard", enableCORS(handleAPI));
	http.HandleFunc("/dashboard/", enableCORS(deleteAPI));
	http.HandleFunc("/dashboard/api/", enableCORS(getChecks));
	
	go startMonitoring();
	
	fmt.Println("Server running on http://localhost:8080");
	 http.ListenAndServe(":8080", nil);

}