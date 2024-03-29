package kerl_test

import (
	. "github.com/iotaledger/iota.go/consts"
	. "github.com/iotaledger/iota.go/kerl"
	. "github.com/iotaledger/iota.go/signing/utils"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kerl", func() {

	Context("hash valid trits", func() {
		DescribeTable("hash computation",
			func(in Trytes, expected Trytes) {
				k := NewKerl()
				Expect(k.Absorb(MustTrytesToTrits(in))).ToNot(HaveOccurred())
				trits, err := k.Squeeze(len(expected) * HashTrinarySize / HashTrytesSize)
				Expect(err).ToNot(HaveOccurred())
				Expect(MustTritsToTrytes(trits)).To(Equal(expected))
			},
			Entry("normal trytes",
				"HHPELNTNJIOKLYDUW9NDULWPHCWFRPTDIUWLYUHQWWJVPAKKGKOAZFJPQJBLNDPALCVXGJLRBFSHATF9C",
				"DMJWZTDJTASXZTHZFXFZXWMNFHRTKWFUPCQJXEBJCLRZOM9LPVJSTCLFLTQTDGMLVUHOVJHBBUYFD9AXX",
			),
			Entry("normal trytes #2",
				"QAUGQZQKRAW9GKEFIBUD9BMJQOABXBTFELCT9GVSZCPTZOSFBSHPQRWJLLWURPXKNAOWCSVWUBNDSWMPW",
				"HOVOHFEPCIGTOFEAZVXAHQRFFRTPQEEKANKFKIHUKSGRICVADWDMBINDYKRCCIWBEOPXXIKMLNSOHEAQZ",
			),
			Entry("normal trytes #3",
				"MWBLYBSRKEKLDHUSRDSDYZRNV9DDCPN9KENGXIYTLDWPJPKBHQBOALSDH9LEJVACJAKJYPCFTJEROARRW",
				"KXBKXQUZBYZFSYSPDPCNILVUSXOEHQWWWFKZPFCQ9ABGIIQBNLSWLPIMV9LYNQDDYUS9L9GNUIYKYAGVZ",
			),
			Entry("input with 243-trits",
				"EMIDYNHBWMBCXVDEFOFWINXTERALUKYYPPHKP9JJFGJEIUY9MUDVNFZHMMWZUYUSWAIOWEVTHNWMHANBH",
				"EJEAOOZYSAWFPZQESYDHZCGYNSTWXUMVJOVDWUNZJXDGWCLUFGIMZRMGCAZGKNPLBRLGUNYWKLJTYEAQX",
			),
			Entry("output with more than 243-trits",
				"9MIDYNHBWMBCXVDEFOFWINXTERALUKYYPPHKP9JJFGJEIUY9MUDVNFZHMMWZUYUSWAIOWEVTHNWMHANBH",
				"G9JYBOMPUXHYHKSNRNMMSSZCSHOFYOYNZRSZMAAYWDYEIMVVOGKPJBVBM9TDPULSFUNMTVXRKFIDOHUXXVYDLFSZYZTWQYTE9SPYYWYTXJYQ9IFGYOLZXWZBKWZN9QOOTBQMWMUBLEWUEEASRHRTNIQWJQNDWRYLCA",
			),
			Entry("input & output with more than 243-trits",
				"G9JYBOMPUXHYHKSNRNMMSSZCSHOFYOYNZRSZMAAYWDYEIMVVOGKPJBVBM9TDPULSFUNMTVXRKFIDOHUXXVYDLFSZYZTWQYTE9SPYYWYTXJYQ9IFGYOLZXWZBKWZN9QOOTBQMWMUBLEWUEEASRHRTNIQWJQNDWRYLCA",
				"LUCKQVACOGBFYSPPVSSOXJEKNSQQRQKPZC9NXFSMQNRQCGGUL9OHVVKBDSKEQEBKXRNUJSRXYVHJTXBPDWQGNSCDCBAIRHAQCOWZEBSNHIJIGPZQITIBJQ9LNTDIBTCQ9EUWKHFLGFUVGGUWJONK9GBCDUIMAYMMQX"),
		)
	})

	Context("hash valid trytes", func() {
		DescribeTable("hash computation",
			func(in Trytes, expected Trytes) {
				k := NewKerl()
				Expect(k.AbsorbTrytes(in)).ToNot(HaveOccurred())
				trytes, err := k.SqueezeTrytes(len(expected) * HashTrinarySize / HashTrytesSize)
				Expect(err).ToNot(HaveOccurred())
				Expect(trytes).To(Equal(expected))
			},
			Entry("normal trytes",
				"HHPELNTNJIOKLYDUW9NDULWPHCWFRPTDIUWLYUHQWWJVPAKKGKOAZFJPQJBLNDPALCVXGJLRBFSHATF9C",
				"DMJWZTDJTASXZTHZFXFZXWMNFHRTKWFUPCQJXEBJCLRZOM9LPVJSTCLFLTQTDGMLVUHOVJHBBUYFD9AXX",
			),
			Entry("normal trytes #2",
				"QAUGQZQKRAW9GKEFIBUD9BMJQOABXBTFELCT9GVSZCPTZOSFBSHPQRWJLLWURPXKNAOWCSVWUBNDSWMPW",
				"HOVOHFEPCIGTOFEAZVXAHQRFFRTPQEEKANKFKIHUKSGRICVADWDMBINDYKRCCIWBEOPXXIKMLNSOHEAQZ",
			),
			Entry("normal trytes #3",
				"MWBLYBSRKEKLDHUSRDSDYZRNV9DDCPN9KENGXIYTLDWPJPKBHQBOALSDH9LEJVACJAKJYPCFTJEROARRW",
				"KXBKXQUZBYZFSYSPDPCNILVUSXOEHQWWWFKZPFCQ9ABGIIQBNLSWLPIMV9LYNQDDYUS9L9GNUIYKYAGVZ",
			),
			Entry("input with 243-trits",
				"EMIDYNHBWMBCXVDEFOFWINXTERALUKYYPPHKP9JJFGJEIUY9MUDVNFZHMMWZUYUSWAIOWEVTHNWMHANBH",
				"EJEAOOZYSAWFPZQESYDHZCGYNSTWXUMVJOVDWUNZJXDGWCLUFGIMZRMGCAZGKNPLBRLGUNYWKLJTYEAQX",
			),
			Entry("output with more than 243-trits",
				"9MIDYNHBWMBCXVDEFOFWINXTERALUKYYPPHKP9JJFGJEIUY9MUDVNFZHMMWZUYUSWAIOWEVTHNWMHANBH",
				"G9JYBOMPUXHYHKSNRNMMSSZCSHOFYOYNZRSZMAAYWDYEIMVVOGKPJBVBM9TDPULSFUNMTVXRKFIDOHUXXVYDLFSZYZTWQYTE9SPYYWYTXJYQ9IFGYOLZXWZBKWZN9QOOTBQMWMUBLEWUEEASRHRTNIQWJQNDWRYLCA",
			),
			Entry("input & output with more than 243-trits",
				"G9JYBOMPUXHYHKSNRNMMSSZCSHOFYOYNZRSZMAAYWDYEIMVVOGKPJBVBM9TDPULSFUNMTVXRKFIDOHUXXVYDLFSZYZTWQYTE9SPYYWYTXJYQ9IFGYOLZXWZBKWZN9QOOTBQMWMUBLEWUEEASRHRTNIQWJQNDWRYLCA",
				"LUCKQVACOGBFYSPPVSSOXJEKNSQQRQKPZC9NXFSMQNRQCGGUL9OHVVKBDSKEQEBKXRNUJSRXYVHJTXBPDWQGNSCDCBAIRHAQCOWZEBSNHIJIGPZQITIBJQ9LNTDIBTCQ9EUWKHFLGFUVGGUWJONK9GBCDUIMAYMMQX"),
		)
	})

	Context("hash invalid trits", func() {

		var k SpongeFunction

		BeforeEach(func() {
			k = NewKerl()
		})
		It("should return an error with empty trits slice", func() {
			Expect(k.Absorb(Trits{})).To(HaveOccurred())
		})

		It("should return an error with invalid trits slice length", func() {
			Expect(k.Absorb(Trits{1, 0, 0, 0, 0, -1})).To(HaveOccurred())
		})
	})

	Context("hash invalid trytes", func() {

		var k SpongeFunction

		BeforeEach(func() {
			k = NewKerl()
		})
		It("should return an error with empty tryte slice", func() {
			Expect(k.AbsorbTrytes("")).To(HaveOccurred())
		})

		It("should return an error with invalid trits slice length", func() {
			Expect(k.AbsorbTrytes("AR")).To(HaveOccurred())
		})
	})

})
