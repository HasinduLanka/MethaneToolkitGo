package main

type ByteSource struct {
	Get func() []byte
}

type BytePipe struct {
	Proc func([]byte) []byte
}

type ByteDestination struct {
	Put func([]byte)
}

func SourcePlusPipe(Source *ByteSource, Pipe *BytePipe) ByteSource {
	return ByteSource{
		Get: func() []byte {
			B := Source.Get()
			return Pipe.Proc(B)
		},
	}
}

func PipePlusPipe(Pipe1 *BytePipe, Pipe2 *BytePipe) BytePipe {
	return BytePipe{
		Proc: func(data []byte) []byte {
			return Pipe2.Proc(Pipe1.Proc(data))
		},
	}
}

func PipePlusDestination(Pipe1 *BytePipe, Pipe2 *ByteDestination) ByteDestination {
	return ByteDestination{
		Put: func(data []byte) {
			Pipe2.Put(Pipe1.Proc(data))
		},
	}
}

func ProcessByteSource(Source *ByteSource, Destination *ByteDestination) {

	data := Source.Get()
	Destination.Put(data)

}

func ProcessByteSources(Sources chan *ByteSource, Destination *ByteDestination) {
	for src := range Sources {
		ProcessByteSource(src, Destination)
	}
}

func NewBytePipe_DecryptAES(key []byte) BytePipe {
	return BytePipe{
		Proc: func(data []byte) []byte {
			return DecryptAES(key, data)
		},
	}
}

func NewBytePipe_EncryptAES(key []byte) BytePipe {
	return BytePipe{
		Proc: func(data []byte) []byte {
			return EncryptAES(key, data)
		},
	}
}

func NewByteDestination_WriteToFile(filename string) ByteDestination {
	return ByteDestination{
		Put: func(data []byte) {
			WriteFile(filename, data)
		},
	}
}

func NewByteDestination_AppendToFile(filename string) ByteDestination {
	return ByteDestination{
		Put: func(data []byte) {
			if len(data) != 0 {
				AppendFile(filename, data)
			}
		},
	}
}

func NewByteDestination_WriteToFiles(FileNames *LineProvider) ByteDestination {
	return ByteDestination{
		Put: func(data []byte) {
			filename, ok := <-FileNames.ch
			if !ok {
				filename = UniqueStamp() + "-nameless"
				Print("File names not enough")
			}
			WriteFile(filename, data)
		},
	}
}
