package OpenPgpJsApi

import (
	"fmt"
)

func ExampleExportPublicKeysRequest_Execute() {
	request := ExportPublicKeysRequest{"1E43F132357B5AD55CECCCC3067D1766157F6495"};
	result, err := request.Execute()
	if err != nil {
		panic(fmt.Sprintf("export request errored: %s", err))
	}

	fmt.Print(result.KeyBlock)

	// Output:
	//-----BEGIN PGP PUBLIC KEY BLOCK-----
	//
	//mI0EWQROJwEEAMtfVOyRBvd1BToqLxBnpbFRoZ3JebXjevOrojC21t9Qcxf7oTSi
	//LevdzvuUQUUF9zoJY0G83UEYLAQHwzaqVYRi9kEfkAX4jYfvuRxCCBNwYtsXVC8b
	//XfQQLVJhVjuhPKtV2AYmsGyJgpVvKjU3HkC7CRr+6PKhugenfaMaXqblABEBAAG0
	//QlRlc3QgS2V5IChvbmx5IGZvciB0ZXN0IHB1cnBvc2VzKSA8Z251cGdob3N0YXBw
	//X3Rlc3RzQGV4YW1wbGUuY29tPoi4BBMBAgAiBQJZBE4nAhsDBgsJCAcDAgYVCAIJ
	//CgsEFgIDAQIeAQIXgAAKCRAGfRdmFX9klfwNBACc2LaE6OFyBiXra405jujKociE
	//TNWYveuIB7p6mGnh+ssoswWPKd02uO5OxQayBbM3WA0mDPe3PtBXwbjFG6QnSv5C
	//eVZejtQvax06uyw48jd1naqz609iNx/buc5NP6rQ50WzmaPk6C3anPd3nICOZufz
	//TuQd0ZILly1xS8bRH7iNBFkETicBBADB7ZHODrmDqJ5mY4ybQI7FN1bTdh24Hpje
	//FcHTRrvQApN4Ttm/IM07cvKWUppQwuzMwvpjxPPloB/oImpr36wDPmDN6lotsQK8
	//W5HHKTEAUoDJoLGXgVuafnetr+q8hfvi/jsuw1GKGU2cJkQdm9Bw7z1ppmlTLprh
	//TbY3s7GsvQARAQABiJ8EGAECAAkFAlkETicCGwwACgkQBn0XZhV/ZJU8cwP/bTIo
	//cM3fr3iWU3bdo1zmXnzSf7kIsrHUTfZqshfJyDIYJgQaTdGav9Uq/Ncwjxlrnw40
	//DuIouEvacGLPXUUnDMPXKBPPRwNvVdrKx1fZVH4jERI16P0Fjq2u2Gisvb3WBJMM
	//lyEL7Mb18KJLTFMGJc3P6nu61b4wLrKccOvOKjc=
	//=xXDW
	//-----END PGP PUBLIC KEY BLOCK-----
}


