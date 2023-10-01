package graph

// var members []*model.Member
// query := r.DB.Model(&model.Member{})/

// 	// Assuming you have a GORM database connection named "db"
// 	err = query.Find(&members).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Return the list of members with populated registrations
// 	return members, nil
// 	}


// 	func (r *queryResolver) GetChurchByID(ctx context.Context, id string) ( *model.Church,  error){
// 		// Extract Admin's ID from the request context (provided by AuthenticationMiddleware).
// 	err := middleware.ExtractCTXinfo(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	church := &model.Church{}

// 	// Use Gorm to find the church by its ID
// 	if err := r.DB.Where("id = ?", id).Preload("SubChurches").First(church).Error; err != nil {
// 		return nil, err
// 	}

// 	return church, nil
// 	}


// 	func (r *queryResolver) GetsubChurchByID(ctx context.Context, id string) ( *model.SubChurch,  error){
// 		// Extract Admin's ID from the request context (provided by AuthenticationMiddleware).
// 	err := middleware.ExtractCTXinfo(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	church := &model.SubChurch{}

// 	// Use Gorm to find the church by its ID
// 	if err := r.DB.Where("id = ?", id).Preload("Members.Registrations").First(church).Error; err != nil {
// 		return nil, err
// 	}

// 	return church, nil
// 	}


// 	func (r *queryResolver) GetAllChurchByID(ctx context.Context) ( []*model.Church,  error){
// 		// Extract Admin's ID from the request context (provided by AuthenticationMiddleware).
// 	// err := middleware.ExtractCTXinfo(ctx)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
	



// var churchs []*model.Church
// if err := r.DB.Preload("SubChurches.Members").Find(&churchs).Error; err != nil {
// 		return nil, err
// 	}

// 	return churchs, nil
// 	}


// 	func (r *queryResolver) GetsubChurchByMainChurchID(ctx context.Context, mainChurchID string) ( []*model.SubChurch,  error){
		



// var subChurch []*model.SubChurch
// if err := r.DB.Where("church_id = ?", mainChurchID).Find(&subChurch).Error; err != nil {
// 		return nil, err
// 	}

// 	return subChurch, nil
// 	}


// 	func (r *queryResolver) MembersBySubChurchID(ctx context.Context, subChurchID string) ( []*model.Member,  error){
// 		// Extract Admin's ID from the request context (provided by AuthenticationMiddleware).
// 	// err := middleware.ExtractCTXinfo(ctx)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
	



// var members []*model.Member
// if err := r.DB.
// 		Where("sub_church_id = ?", subChurchID). // Specify the field to compare against subChurchID
// 		Preload("Registrations").
// 		Preload("Church").
// 		Find(&members).Error; err != nil {
// 		fmt.Println("Error preloading associations:", err)
// 		return nil, err
// 	}

// 	return members, nil
// 	}


// 	func (r *queryResolver) GetAllsubChurch(ctx context.Context) ( []*model.SubChurch,  error){
// 		// Extract Admin's ID from the request context (provided by AuthenticationMiddleware).
// 	err := middleware.ExtractCTXinfo(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

	



// var subs []*model.SubChurch
// if err := r.DB.Preload("Members").Find(&subs).Error; err != nil {
// 		return nil, err
// 	}

// 	return subs, nil
// 	}


// 	func (r *queryResolver) GetsubChurch(ctx context.Context, id string) ( *model.SubChurch,  error){
// 		// Extract Admin's ID from the request context (provided by AuthenticationMiddleware).
// 	err := middleware.ExtractCTXinfo(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
	



// var subChurch model.SubChurch
// if err := r.DB.Preload("Members.Registrations").First(&subChurch, "id = ?", id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &subChurch, nil
// 	}


// 	func (r *queryResolver) CurrentWeekRegistrations(ctx context.Context) ( []*model.Registration,  error){
// 		// Extract Admin's ID from the request context (provided by AuthenticationMiddleware).
// 	err := middleware.ExtractCTXinfo(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

	



// var registrations []*model.Registration
// if err := r.DB.Preload("Member").Find(&registrations).Error; err != nil {
// 		return nil, err
// 	}
// 	currentWeekNumber := schemas.GetWeekNumber(time.Now()) // Get the current week number

	



// var currentWeekRegistrations []*model.Registration
// for _, reg := range registrations {
// 		if schemas.GetWeekNumber(reg.CreatedAt) == currentWeekNumber {
// 			currentWeekRegistrations = append(currentWeekRegistrations, reg)
// 		}
// 	}
// 	return currentWeekRegistrations, nil
// 	}


// 	func (r *queryResolver) GetRegistrations(ctx context.Context) ( []*model.WeeklyResults,  error){
// 		// Extract Admin's ID from the request context (provided by AuthenticationMiddleware).
// 	err := middleware.ExtractCTXinfo(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
	



// var registrations []*model.Registration
// if err := r.DB.Find(&registrations).Error; err != nil {
// 		return nil, err
// 	}

// 	registrationsByWeek := make(map[int][]*model.Registration)

// 	for _, reg := range registrations {
// 		if reg.Absence { // Filter only registrations with absence true
// 			weekNumber := schemas.GetWeekNumber(reg.CreatedAt)
// 			registrationsByWeek[weekNumber] = append(registrationsByWeek[weekNumber], reg)
// 		}
// 	}

	



// var monthlyResults []*model.WeeklyResults
// for _, regs := range registrationsByWeek {
// 		year, month, day := regs[0].CreatedAt.Date()
// 		weekOfMonth := (day-1)/7 + 1
// 		monthStr := month.String() + " " + strconv.Itoa(year) // Get the month string
// 		monthlyResults = append(monthlyResults, &model.WeeklyResults{
// 			Date: &model.DateInfo{
// 				Month:       &monthStr, // Use a pointer to the month string
// 				WeekOfMonth: &weekOfMonth,
// 			},
// 			Registrations: regs,
// 		})
// 	}

// 	return monthlyResults, nil
// 	}


// 	func (r *registrationResolver) ID(ctx context.Context, obj *model.Registration) ( string,  error){
// 		id := obj.ID.String()
// 	return id, nil
// 	}


// 	func (r *subChurchResolver) ID(ctx context.Context, obj *model.SubChurch) ( string,  error){
// 		id := obj.ID.String()
// 	return id, nil
// 	}




// 	func (r *Resolver) Church() ChurchResolver { return &churchResolver{r} }

// 	func (r *Resolver) Member() MemberResolver { return &memberResolver{r} }

// 	func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// 	func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// 	func (r *Resolver) Registration() RegistrationResolver { return &registrationResolver{r} }

// 	func (r *Resolver) SubChurch() SubChurchResolver { return &subChurchResolver{r} }



// func (r *queryResolver) MembersBySubChurch(ctx context.Context, subChurchID string) ([]*model.Member, error) {
// 	panic(fmt.Errorf("not implemented: MembersBySubChurch - membersBySubChurch"))
// }

