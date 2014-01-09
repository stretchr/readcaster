# ReadCaster

ReadCaster is a go package (for Golang) for elegantly broadcasting the data from one io.Reader source to many `io.Readers` in a memory efficient way.

  * Jump straight into the [API Documentation](http://godoc.org/github.com/stretchr/readcaster)

### Problem

We needed to stream content from an `io.Reader` to many external programs and didn't want to keep multiple copies of the data in memory.  This solution allows us to control how much memory is used at one time.

### How it works

ReadCaster reads a buffer load of data from the source, and channels it to each of the readers.  If one reader is slower than the others, the other readers will be blocked to give that reader time to catch-up.  This allows ReadCaster to ensure it doesn't use an ever growing amount of memory.

## Usage

    // make a source (can be any io.Reader)
    source := strings.NewReader("Hello from Stretchr!")

    // make a caster for this source
    caster := readcaster.New(source)

    // make three readers that need to read the same source
    reader1 := caster.NewReader()
    reader2 := caster.NewReader()
    reader3 := caster.NewReader()

    // use a wait group so we can wait for all readers to
    // be finished
    var waiter sync.WaitGroup
    waiter.Add(3)

    // we're going to read three lots of the data
    var bytes1, bytes2, bytes3 []byte

    // trigger goroutines to read from the three
    // readers
    go func(){
      bytes1, _ = ioutil.ReadAll(reader1)
      waiter.Done()
    }()
    go func(){
      bytes2, _ = ioutil.ReadAll(reader2)
      waiter.Done()
    }()
    go func(){
      bytes3, _ = ioutil.ReadAll(reader3)
      waiter.Done()
    }()

    // Now we have three copies of the content
    // and we have only read from one source.
    //
    // And the memory footprint is predictable and controlable.
    fmt.Sprintf("I only used %dK memory", caster.ApproxMemoryUse())

### Controlling memory

To be specific about the amount of memory you plan to use, you may use the `NewSize` method.

    readcaster.NewSize(source, bufferSize, backlogSize)

  * `source` - the `io.Reader` to read from
  * `bufferSize` - the size (in bytes) of each buffer (default is `4096` or 4kb)
  * `backlogSize` - the number of buffers that will be queued up for each reader before the reader gets blocked if other readers are being slow (default is `10`)

Calculating the `ApproxMemoryUse()` is as simple as:

    bytesUsed := bufferSize * backlogSize * numberOfReaders

  * `numberOfReaders` is the number of readers that were generated by calls to `caster.NewReader`.

So for three readers at the default settings, it would be:

    4096 * 10 * 3 = 120kb
