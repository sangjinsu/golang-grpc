# gRPC Golang Master Class

section 10 CRUD API with MongoDB

### 프로그램 실행 팁

- 버그나 에러 발생시 파일 이름과 줄 번호를 표시

  ```go
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  ```

- crtl + c 신호 대기 및 신호시 종료

  ```go
  	go func() {
  		fmt.Println("Starting Server...")
  		if serveErr := s.Serve(lis); serveErr != nil {
  			log.Fatalf("Failed to serve: %v", serveErr)
  		}
  	}()
  
  	// Wait for Control C to exit
  	// ctrl + c 대기
  	ch := make(chan os.Signal, 1)
  	signal.Notify(ch, os.Interrupt)
  
  	// Block until a signal is received
  	<-ch
  	fmt.Println("Stopping the server")
  	s.Stop()
  	fmt.Println("Closing the listener")
  	closeErr := lis.Close()
  	if closeErr != nil {
  		log.Fatalf("Failed to close: %v", closeErr)
  	}
  	fmt.Println("End of Program")
  ```

  

