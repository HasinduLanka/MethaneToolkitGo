package methane

// Keep key nil for unencrypted
func URIListParser(PathPrefixForURIs string, URIProvider *LineProvider, IsEncrypted bool, key []byte, OutFileName string) {
	lp := NewLinePipe_FilterOutShellStyleComments().Proc(URIProvider)
	lp = NewLinePipe_TrimAndFilterOutEmpty().Proc(lp)
	lp = NewLinePipe_Prefix(PathPrefixForURIs).Proc(lp)

	var BD ByteDestination

	if IsEncrypted {
		Decrypt := NewBytePipe_DecryptAES(key)
		Merge := NewByteDestination_AppendToFile(WSRoot + OutFileName)

		BD = PipePlusDestination(&Decrypt, &Merge)

	} else {
		BD = NewByteDestination_AppendToFile(WSRoot + OutFileName)
	}

	ByteSrcs := URIsToByteSources(lp)

	ProcessByteSources(ByteSrcs, &BD)

}

// TODO : User interface for Parsing M3U files

// func RunM3U() {

// 	var filename string
// 	if len(os.Args) > 1 {
// 		filename = os.Args[1]
// 	} else {
// 		filename = Prompt("Enter file name or URL : ")
// 		if len(filename) == 0 {
// 			filename = wsroot + "playlist.m3u8"
// 		}
// 	}

// 	Print("Reading File " + filename)
// 	file, lurierr := LoadURIToString(filename)
// 	CheckError(lurierr)

// 	keyMethod, keyURI := m3u_ParseKey(file)
// 	var key []byte

// 	if len(keyMethod) == 0 {
// 		r := PromptOptions("Key info not found", map[string]string{"p": "I'll provide key data", "n": "Do not use encryption"})
// 		switch r {
// 		case "p":
// 			keyMethod, key = m3u_PromptKeyData()
// 		case "n":
// 			keyMethod = ""
// 			key = nil
// 			keyURI = ""
// 		}

// 	} else {
// 		r := PromptOptions("Key found Method:"+keyMethod+" URI:"+keyURI, map[string]string{"r": "Read this key (Default)", "p": "I'll provide key data", "n": "Do not use encryption"})
// 		switch r {
// 		case "r":
// 			keyMethod, key = m3u_GetKeyData(keyURI)
// 			if len(keyMethod) == 0 {
// 				Print("Reading Key from " + keyURI + " failed")
// 				keyMethod, key = m3u_PromptKeyData()
// 			}
// 		case "p":
// 			keyMethod, key = m3u_PromptKeyData()
// 		case "n":
// 			keyMethod = ""
// 			key = nil
// 			keyURI = ""
// 		}
// 	}

// 	UseAES := (len(keyMethod) != 0)

// 	if UseAES && !TestEncryptionKey(key) {
// 		panic(errors.New("encryption key failed the simple test"))
// 	} else {
// 		Print("Key is looking good. Method:" + keyMethod)
// 	}

// 	// tstf, _ := LoadURI(wsroot + "chunk0.ts")
// 	// decrf := DecryptAES(key, tstf)
// 	// WriteFile(wsroot+"dchunk0.ts", decrf)

// 	var outFileName string

// 	if NoConsole {
// 		outputDir := wsroot + "output/"
// 		MakeDir(outputDir)
// 		outFileName = outputDir + "out.mkv"
// 	} else {
// 		outFileName = Prompt("Enter output file name : ")
// 	}

// 	os.Rename(outFileName, outFileName+".old")

// 	DecrypAndMerge(strings.Split(file, "\n"), outFileName, key)

// }

// func DecrypAndMerge(URIList []string, outFileName string, key []byte) error {
// 	for _, uri := range URIList {

// 		if strings.HasPrefix(uri, "#") {
// 			continue
// 		}

// 		Print("Loading " + uri)
// 		IB, err := LoadURI(uri)
// 		if err != nil {
// 			return err
// 		}

// 		OB := DecryptAES(key, IB)
// 		AppendFile(outFileName, OB)

// 	}

// 	return nil

// }

// func m3u_GetKeyData(uri string) (string, []byte) {
// 	B, err := LoadURI(uri)
// 	if err != nil {
// 		return "", nil
// 	} else {
// 		if (len(B) % 8) == 0 {
// 			return "AES-" + strconv.Itoa(len(B)*8), B
// 		} else {
// 			Print("Key is not in the correct length. Keys should be 16, 24, 32 long but this is " + strconv.Itoa(len(B)))
// 			return "", nil
// 		}
// 	}
// }

// func m3u_PromptKeyData() (string, []byte) {
// 	var r string

// 	if NoConsole {
// 		r = wsroot + "video.key"
// 	} else {
// 		r = Prompt("Enter key file name or URL : ")
// 	}

// 	M, B := m3u_GetKeyData(r)

// 	if len(M) == 0 {
// 		Print("Sorry, I didn't get that")
// 		return m3u_PromptKeyData()
// 	} else {
// 		Print("Key loaded from " + r)
// 		return M, B
// 	}

// }

// var regex_m3u_key *regexp.Regexp = regexp.MustCompile("#EXT-X-KEY:METHOD=(.*),URI=(.*)")

// // Returns : (Method string, Key string)
// func m3u_ParseKey(c string) (string, string) {
// 	re := regex_m3u_key.FindStringSubmatch(c)
// 	if re == nil || len(re) != 3 {
// 		return "", ""
// 	}

// 	return re[1], re[2]

// }
