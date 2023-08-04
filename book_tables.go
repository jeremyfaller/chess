package main

import (
	"math/rand"
)

func getBook(b *Board, r *rand.Rand) (Move, bool) {
	if v, ok := book[b.ZHash()]; ok {
		idx := r.Intn(v.total)
		for _, test := range v.moves {
			idx -= test.count
			if idx < 0 {
				move := Move{from: test.from, to: test.to, p: b.at(test.from)}
				if !b.isLegalMove(&move) {
					panic("move should be legal")
				}
				return move, true
			}
		}
		panic("shouldn't be reachable")
	}
	return Move{}, false
}

type shortMove struct {
	count    int
	from, to Coord
}

type Book struct {
	total int
	moves []shortMove
}

var book = map[Hash]Book{
	7212155634298086293: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 243109},
			{from: 11, to: 27, count: 146627},
			{from: 6, to: 21, count: 33009},
			{from: 10, to: 26, count: 22211},
			{from: 13, to: 29, count: 4982},
		},
		total: 449938,
	},
	13834452002831339266: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 827},
			{from: 62, to: 45, count: 569},
			{from: 50, to: 34, count: 235},
			{from: 52, to: 36, count: 222},
			{from: 54, to: 46, count: 156},
			{from: 52, to: 44, count: 146},
			{from: 51, to: 43, count: 88},
			{from: 50, to: 42, count: 88},
			{from: 57, to: 42, count: 71},
			{from: 53, to: 37, count: 46},
			{from: 49, to: 41, count: 30},
		},
		total: 2478,
	},
	6668294972361061752: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 46},
		},
		total: 46,
	},
	4367353172880353314: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 42},
			{from: 11, to: 27, count: 22},
		},
		total: 64,
	},
	2666486091856982072: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 64},
			{from: 6, to: 21, count: 28},
			{from: 11, to: 27, count: 24},
		},
		total: 116,
	},
	10912508667191605942: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 92},
			{from: 11, to: 27, count: 23},
		},
		total: 115,
	},
	7139621755721538125: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 30},
			{from: 11, to: 27, count: 29},
		},
		total: 59,
	},
	12873515539516911254: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 382},
			{from: 11, to: 27, count: 79},
			{from: 6, to: 21, count: 55},
		},
		total: 516,
	},
	3139134565908222192: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 89},
			{from: 54, to: 46, count: 78},
			{from: 52, to: 44, count: 26},
		},
		total: 193,
	},
	224245799301391684: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 60},
		},
		total: 60,
	},
	14025309784341266392: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 27},
		},
		total: 27,
	},
	17536505249008392295: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 30},
		},
		total: 30,
	},
	18239896499613277683: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 11422},
			{from: 51, to: 35, count: 9052},
			{from: 50, to: 34, count: 3636},
			{from: 51, to: 43, count: 1799},
			{from: 54, to: 46, count: 1791},
			{from: 52, to: 44, count: 1763},
			{from: 50, to: 42, count: 1113},
			{from: 57, to: 42, count: 869},
			{from: 53, to: 37, count: 742},
			{from: 49, to: 41, count: 434},
		},
		total: 32621,
	},
	6830313020568341943: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 185},
			{from: 11, to: 27, count: 75},
			{from: 12, to: 28, count: 63},
			{from: 10, to: 26, count: 56},
			{from: 9, to: 17, count: 28},
		},
		total: 407,
	},
	3508395161547072251: Book{
		moves: []shortMove{
			{from: 58, to: 49, count: 26},
		},
		total: 26,
	},
	17265528953300545948: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 22},
		},
		total: 22,
	},
	15558465085432270846: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 40},
		},
		total: 40,
	},
	15758367162317041258: Book{
		moves: []shortMove{
			{from: 12, to: 20, count: 30},
			{from: 14, to: 22, count: 26},
		},
		total: 56,
	},
	10478139460285737812: Book{
		moves: []shortMove{
			{from: 52, to: 44, count: 25},
		},
		total: 25,
	},
	834659513188881518: Book{
		moves: []shortMove{
			{from: 5, to: 12, count: 37},
			{from: 10, to: 26, count: 20},
		},
		total: 57,
	},
	13677636687619183392: Book{
		moves: []shortMove{
			{from: 61, to: 52, count: 21},
		},
		total: 21,
	},
	4461484135476757359: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 29},
		},
		total: 29,
	},
	404161395204259085: Book{
		moves: []shortMove{
			{from: 58, to: 49, count: 31},
		},
		total: 31,
	},
	18225725712847734118: Book{
		moves: []shortMove{
			{from: 58, to: 49, count: 22},
		},
		total: 22,
	},
	10157261298205427210: Book{
		moves: []shortMove{
			{from: 58, to: 49, count: 180},
		},
		total: 180,
	},
	1587885421407779120: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 65},
		},
		total: 65,
	},
	10699058090883083852: Book{
		moves: []shortMove{
			{from: 58, to: 49, count: 63},
		},
		total: 63,
	},
	5427287738000866603: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 120},
		},
		total: 120,
	},
	34215131378681226: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 46},
			{from: 53, to: 37, count: 42},
		},
		total: 88,
	},
	946807717401551902: Book{
		moves: []shortMove{
			{from: 10, to: 26, count: 102},
			{from: 11, to: 19, count: 75},
			{from: 11, to: 27, count: 50},
		},
		total: 227,
	},
	12612139314359102671: Book{
		moves: []shortMove{
			{from: 61, to: 52, count: 48},
		},
		total: 48,
	},
	12072103605345087040: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 21},
		},
		total: 21,
	},
	6043547581280064877: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 179},
		},
		total: 179,
	},
	15108735916296651281: Book{
		moves: []shortMove{
			{from: 52, to: 44, count: 40},
			{from: 62, to: 45, count: 39},
			{from: 54, to: 46, count: 35},
		},
		total: 114,
	},
	17951875160985040805: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 32},
		},
		total: 32,
	},
	12829005271748644612: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 32},
		},
		total: 32,
	},
	11317136528029037405: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 20},
		},
		total: 20,
	},
	956403182426560350: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 31},
		},
		total: 31,
	},
	15922249594438165381: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 324},
		},
		total: 324,
	},
	10943508153249225508: Book{
		moves: []shortMove{
			{from: 52, to: 44, count: 194},
			{from: 50, to: 34, count: 59},
			{from: 54, to: 46, count: 51},
		},
		total: 304,
	},
	9344940315755145118: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 303},
		},
		total: 303,
	},
	278124154402451682: Book{
		moves: []shortMove{
			{from: 58, to: 49, count: 300},
		},
		total: 300,
	},
	2186061131103709065: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 346},
			{from: 10, to: 26, count: 328},
			{from: 11, to: 27, count: 250},
			{from: 12, to: 28, count: 61},
			{from: 9, to: 17, count: 57},
			{from: 11, to: 19, count: 20},
		},
		total: 1062,
	},
	8082840153283643589: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 57},
		},
		total: 57,
	},
	13579090082173038424: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 33},
		},
		total: 33,
	},
	6165234546481193655: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 42},
		},
		total: 42,
	},
	14729446791029304372: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 307},
			{from: 51, to: 43, count: 29},
		},
		total: 336,
	},
	5696962883902032660: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 49},
		},
		total: 49,
	},
	14806797955201375336: Book{
		moves: []shortMove{
			{from: 57, to: 51, count: 20},
		},
		total: 20,
	},
	13845862591266979324: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 37},
		},
		total: 37,
	},
	13921667893093060000: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 55},
		},
		total: 55,
	},
	4845960014123703004: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 40},
		},
		total: 40,
	},
	9099256295467678419: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 763},
			{from: 14, to: 22, count: 465},
			{from: 10, to: 26, count: 290},
			{from: 12, to: 28, count: 96},
			{from: 9, to: 17, count: 72},
			{from: 11, to: 19, count: 51},
			{from: 12, to: 20, count: 24},
		},
		total: 1761,
	},
	1170225804422124959: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 27},
			{from: 62, to: 45, count: 24},
		},
		total: 51,
	},
	2118742182357103627: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 73},
		},
		total: 73,
	},
	2672594188682824297: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 38},
			{from: 52, to: 36, count: 36},
		},
		total: 74,
	},
	15880783122503823874: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 42},
			{from: 52, to: 36, count: 27},
		},
		total: 69,
	},
	15076313213518713750: Book{
		moves: []shortMove{
			{from: 1, to: 11, count: 23},
			{from: 14, to: 22, count: 22},
		},
		total: 45,
	},
	12427964308317434222: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 143},
			{from: 62, to: 45, count: 138},
			{from: 54, to: 46, count: 86},
			{from: 50, to: 42, count: 26},
			{from: 53, to: 37, count: 25},
			{from: 58, to: 30, count: 21},
		},
		total: 439,
	},
	9580321535421957338: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 86},
		},
		total: 86,
	},
	470589228595221414: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 83},
		},
		total: 83,
	},
	2990563377236848639: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 316},
			{from: 11, to: 27, count: 41},
		},
		total: 357,
	},
	7109084321425938270: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 92},
			{from: 57, to: 51, count: 59},
			{from: 52, to: 36, count: 56},
			{from: 50, to: 42, count: 43},
			{from: 57, to: 42, count: 28},
		},
		total: 278,
	},
	12954752049037370777: Book{
		moves: []shortMove{
			{from: 11, to: 19, count: 28},
			{from: 11, to: 27, count: 22},
		},
		total: 50,
	},
	8066712635454797514: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 214},
			{from: 11, to: 19, count: 159},
			{from: 10, to: 26, count: 106},
		},
		total: 479,
	},
	14832136015451110939: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 145},
		},
		total: 145,
	},
	8169821528847968407: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 221},
			{from: 1, to: 11, count: 167},
			{from: 10, to: 26, count: 88},
			{from: 10, to: 18, count: 58},
		},
		total: 534,
	},
	5935357386665155160: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 64},
			{from: 50, to: 34, count: 48},
			{from: 57, to: 42, count: 30},
			{from: 57, to: 51, count: 29},
		},
		total: 171,
	},
	10283531027212144407: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 28},
		},
		total: 28,
	},
	2550345163500147612: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 29},
			{from: 50, to: 34, count: 20},
			{from: 50, to: 42, count: 20},
		},
		total: 69,
	},
	855099580379787826: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 125},
		},
		total: 125,
	},
	4680923931715406483: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 121},
		},
		total: 121,
	},
	12247822798765095651: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 31},
		},
		total: 31,
	},
	9727284904867612346: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 52},
		},
		total: 52,
	},
	8111034973309674529: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 22},
		},
		total: 22,
	},
	10768388569329624080: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 21},
		},
		total: 21,
	},
	11611068534874575098: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 380},
			{from: 11, to: 27, count: 34},
		},
		total: 414,
	},
	2544868681062074246: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 211},
			{from: 52, to: 36, count: 83},
			{from: 57, to: 51, count: 36},
			{from: 50, to: 42, count: 33},
			{from: 57, to: 42, count: 23},
		},
		total: 386,
	},
	17451796322423043393: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 35},
		},
		total: 35,
	},
	13337739329353846240: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 29},
		},
		total: 29,
	},
	7484343830100353737: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 593},
			{from: 10, to: 26, count: 519},
			{from: 11, to: 27, count: 359},
			{from: 9, to: 17, count: 101},
			{from: 12, to: 28, count: 87},
			{from: 12, to: 20, count: 36},
			{from: 11, to: 19, count: 35},
		},
		total: 1730,
	},
	695749242735644037: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 62},
		},
		total: 62,
	},
	359723514037873681: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 84},
		},
		total: 84,
	},
	4362932689970680435: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 72},
			{from: 61, to: 52, count: 50},
		},
		total: 122,
	},
	14265528365422815768: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 38},
		},
		total: 38,
	},
	14457509825000328076: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 20},
			{from: 1, to: 11, count: 20},
		},
		total: 40,
	},
	3171827355698554871: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 53},
		},
		total: 53,
	},
	2404547050735863395: Book{
		moves: []shortMove{
			{from: 9, to: 17, count: 34},
			{from: 11, to: 27, count: 20},
		},
		total: 54,
	},
	13106119068265261428: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 380},
			{from: 49, to: 41, count: 64},
			{from: 62, to: 45, count: 49},
			{from: 53, to: 37, count: 44},
		},
		total: 537,
	},
	7569030618696547387: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 20},
		},
		total: 20,
	},
	13310524193007782112: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 460},
		},
		total: 460,
	},
	4236628301145841564: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 338},
			{from: 61, to: 52, count: 63},
			{from: 50, to: 34, count: 36},
		},
		total: 437,
	},
	3622140227679198483: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 60},
		},
		total: 60,
	},
	8749502672033441202: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 51},
		},
		total: 51,
	},
	14257031002429698878: Book{
		moves: []shortMove{
			{from: 11, to: 19, count: 21},
		},
		total: 21,
	},
	15397241043282864199: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 572},
			{from: 11, to: 27, count: 524},
			{from: 10, to: 26, count: 357},
			{from: 12, to: 28, count: 146},
			{from: 15, to: 31, count: 54},
			{from: 9, to: 17, count: 50},
			{from: 12, to: 20, count: 36},
			{from: 11, to: 19, count: 25},
			{from: 10, to: 18, count: 20},
		},
		total: 1784,
	},
	13535439392184026891: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 51},
		},
		total: 51,
	},
	10723292320968538962: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 32},
		},
		total: 32,
	},
	13154372983497762463: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 303},
		},
		total: 303,
	},
	10303116660681234685: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 347},
		},
		total: 347,
	},
	11671999718785148068: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 90},
			{from: 12, to: 20, count: 84},
			{from: 10, to: 26, count: 61},
			{from: 12, to: 28, count: 58},
			{from: 11, to: 27, count: 29},
		},
		total: 322,
	},
	16871160149369561498: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 210},
			{from: 51, to: 43, count: 21},
		},
		total: 231,
	},
	7585269701562287802: Book{
		moves: []shortMove{
			{from: 5, to: 12, count: 21},
		},
		total: 21,
	},
	6238957450394466070: Book{
		moves: []shortMove{
			{from: 5, to: 12, count: 98},
			{from: 11, to: 27, count: 63},
			{from: 10, to: 26, count: 44},
		},
		total: 205,
	},
	16199675187773026392: Book{
		moves: []shortMove{
			{from: 51, to: 43, count: 55},
		},
		total: 55,
	},
	7185179427358824312: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 26},
		},
		total: 26,
	},
	8341993837275508505: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 235},
		},
		total: 235,
	},
	14948550284138396053: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 214},
		},
		total: 214,
	},
	5548459536378031849: Book{
		moves: []shortMove{
			{from: 51, to: 43, count: 117},
			{from: 51, to: 35, count: 87},
			{from: 50, to: 34, count: 33},
		},
		total: 237,
	},
	14872639497361787337: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 47},
			{from: 4, to: 6, count: 45},
		},
		total: 92,
	},
	9605665466771796328: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 31},
		},
		total: 31,
	},
	18254386348903556001: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 42},
		},
		total: 42,
	},
	15155695008618078200: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 39},
		},
		total: 39,
	},
	7255943715649826082: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 86},
		},
		total: 86,
	},
	5455277215388045691: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 86},
		},
		total: 86,
	},
	8613433761217716374: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 64},
		},
		total: 64,
	},
	6381129781672975567: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 39},
		},
		total: 39,
	},
	8850455900677816578: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 48},
		},
		total: 48,
	},
	6479707630331884493: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 25},
		},
		total: 25,
	},
	8569162784718475156: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 20},
		},
		total: 20,
	},
	12172333049811389119: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 52},
		},
		total: 52,
	},
	9794084053785573094: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 43},
		},
		total: 43,
	},
	403193436182028698: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 155},
			{from: 51, to: 43, count: 21},
		},
		total: 176,
	},
	13343939408685561622: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 84},
			{from: 1, to: 11, count: 32},
			{from: 12, to: 28, count: 30},
		},
		total: 146,
	},
	11101543835836171737: Book{
		moves: []shortMove{
			{from: 51, to: 43, count: 22},
		},
		total: 22,
	},
	17439949018381881271: Book{
		moves: []shortMove{
			{from: 51, to: 43, count: 230},
			{from: 51, to: 35, count: 188},
			{from: 50, to: 34, count: 60},
		},
		total: 478,
	},
	11417995068795219321: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 98},
		},
		total: 98,
	},
	12786737360574796064: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 70},
		},
		total: 70,
	},
	10605674637962972397: Book{
		moves: []shortMove{
			{from: 9, to: 17, count: 39},
			{from: 11, to: 27, count: 35},
		},
		total: 74,
	},
	545516043969799162: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 556},
		},
		total: 556,
	},
	2924327196257230755: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 511},
			{from: 11, to: 27, count: 33},
		},
		total: 544,
	},
	12314725747242551519: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 195},
			{from: 51, to: 43, count: 166},
			{from: 50, to: 34, count: 106},
			{from: 52, to: 36, count: 69},
			{from: 51, to: 35, count: 61},
			{from: 50, to: 42, count: 45},
		},
		total: 642,
	},
	5301268839970992805: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 56},
		},
		total: 56,
	},
	196457329255483908: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 37},
		},
		total: 37,
	},
	8506843633258143120: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 23},
		},
		total: 23,
	},
	12086679756820162891: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 3368},
			{from: 11, to: 19, count: 137},
			{from: 11, to: 27, count: 121},
			{from: 10, to: 26, count: 101},
			{from: 10, to: 18, count: 46},
		},
		total: 3773,
	},
	17047413462871921130: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 2838},
			{from: 51, to: 43, count: 279},
			{from: 51, to: 35, count: 154},
			{from: 50, to: 34, count: 52},
		},
		total: 3323,
	},
	5774679337674377062: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 989},
			{from: 10, to: 26, count: 871},
			{from: 11, to: 19, count: 808},
			{from: 9, to: 17, count: 64},
			{from: 10, to: 18, count: 46},
		},
		total: 2778,
	},
	154921531842182253: Book{
		moves: []shortMove{
			{from: 51, to: 43, count: 21},
		},
		total: 21,
	},
	17706235359885922880: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 37},
		},
		total: 37,
	},
	5264271757369001164: Book{
		moves: []shortMove{
			{from: 4, to: 6, count: 20},
		},
		total: 20,
	},
	776938087158432366: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 1820},
			{from: 9, to: 17, count: 89},
			{from: 9, to: 25, count: 49},
			{from: 10, to: 26, count: 32},
		},
		total: 1990,
	},
	9852680641951175954: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 1827},
			{from: 51, to: 35, count: 25},
		},
		total: 1852,
	},
	2434743011801680060: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 382},
			{from: 14, to: 22, count: 184},
			{from: 10, to: 26, count: 139},
			{from: 12, to: 28, count: 71},
			{from: 9, to: 17, count: 34},
			{from: 11, to: 19, count: 25},
		},
		total: 835,
	},
	5743081910809667568: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 21},
		},
		total: 21,
	},
	17579241544154989313: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 102},
			{from: 51, to: 35, count: 40},
		},
		total: 142,
	},
	18351097466506551957: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 60},
		},
		total: 60,
	},
	8998135085377222121: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 40},
		},
		total: 40,
	},
	17328465824213859431: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 3934},
			{from: 10, to: 26, count: 3883},
			{from: 11, to: 27, count: 2186},
			{from: 9, to: 17, count: 619},
			{from: 11, to: 19, count: 243},
			{from: 12, to: 20, count: 208},
			{from: 1, to: 18, count: 182},
		},
		total: 11255,
	},
	11385491920504918827: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 316},
			{from: 51, to: 35, count: 131},
			{from: 52, to: 44, count: 59},
			{from: 51, to: 43, count: 48},
			{from: 49, to: 41, count: 30},
			{from: 50, to: 34, count: 25},
		},
		total: 609,
	},
	11709007013741512556: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 24},
		},
		total: 24,
	},
	5935742811181610166: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 131},
			{from: 51, to: 35, count: 124},
			{from: 52, to: 44, count: 32},
			{from: 57, to: 42, count: 23},
		},
		total: 310,
	},
	13520422876812982617: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 191},
			{from: 51, to: 35, count: 108},
			{from: 52, to: 44, count: 60},
			{from: 50, to: 34, count: 30},
			{from: 51, to: 43, count: 30},
		},
		total: 419,
	},
	4090425864158830201: Book{
		moves: []shortMove{
			{from: 9, to: 17, count: 25},
			{from: 11, to: 27, count: 21},
		},
		total: 46,
	},
	2469345947548350426: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 2023},
			{from: 51, to: 35, count: 785},
			{from: 52, to: 44, count: 379},
			{from: 49, to: 41, count: 303},
			{from: 50, to: 34, count: 145},
			{from: 51, to: 43, count: 141},
			{from: 49, to: 33, count: 80},
			{from: 57, to: 42, count: 66},
			{from: 50, to: 42, count: 42},
		},
		total: 3964,
	},
	1446424639641465563: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 145},
			{from: 52, to: 36, count: 101},
			{from: 62, to: 45, count: 56},
			{from: 54, to: 46, count: 42},
			{from: 50, to: 34, count: 33},
		},
		total: 377,
	},
	1821857140141708111: Book{
		moves: []shortMove{
			{from: 1, to: 18, count: 28},
		},
		total: 28,
	},
	751201845693549785: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 1172},
			{from: 51, to: 35, count: 995},
			{from: 62, to: 45, count: 879},
			{from: 50, to: 34, count: 227},
			{from: 52, to: 44, count: 186},
			{from: 51, to: 43, count: 185},
			{from: 49, to: 41, count: 112},
			{from: 54, to: 46, count: 100},
			{from: 50, to: 42, count: 97},
			{from: 48, to: 32, count: 61},
			{from: 57, to: 42, count: 44},
		},
		total: 4058,
	},
	12215955319985675421: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 109},
		},
		total: 109,
	},
	10517721698625348351: Book{
		moves: []shortMove{
			{from: 58, to: 49, count: 108},
		},
		total: 108,
	},
	5682801642865324440: Book{
		moves: []shortMove{
			{from: 12, to: 20, count: 48},
			{from: 6, to: 21, count: 31},
		},
		total: 79,
	},
	361553369572947110: Book{
		moves: []shortMove{
			{from: 52, to: 44, count: 26},
			{from: 62, to: 45, count: 22},
		},
		total: 48,
	},
	11527214942712105884: Book{
		moves: []shortMove{
			{from: 6, to: 21, count: 20},
		},
		total: 20,
	},
	490747803628369402: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 25},
		},
		total: 25,
	},
	16797190927209739939: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 91},
		},
		total: 91,
	},
	15089911623788349633: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 81},
		},
		total: 81,
	},
	9874953928002497529: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 179},
		},
		total: 179,
	},
	12788917599045850523: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 80},
			{from: 62, to: 45, count: 43},
			{from: 57, to: 51, count: 22},
		},
		total: 145,
	},
	13556237494882333711: Book{
		moves: []shortMove{
			{from: 12, to: 20, count: 72},
			{from: 14, to: 22, count: 28},
			{from: 9, to: 45, count: 25},
			{from: 6, to: 21, count: 25},
		},
		total: 150,
	},
	17868116248245121329: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 26},
		},
		total: 26,
	},
	10416918487581406179: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 178},
		},
		total: 178,
	},
	12169265673147272577: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 105},
			{from: 62, to: 45, count: 50},
		},
		total: 155,
	},
	11942376632948894741: Book{
		moves: []shortMove{
			{from: 12, to: 20, count: 152},
			{from: 6, to: 21, count: 75},
			{from: 14, to: 22, count: 28},
		},
		total: 255,
	},
	17177568343035210027: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 72},
			{from: 61, to: 52, count: 43},
		},
		total: 115,
	},
	15801291677559092957: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 24},
		},
		total: 24,
	},
	2513028407064344941: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 98},
		},
		total: 98,
	},
	1914501045358708495: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 91},
		},
		total: 91,
	},
	1714638286601594523: Book{
		moves: []shortMove{
			{from: 9, to: 45, count: 224},
			{from: 12, to: 20, count: 101},
			{from: 12, to: 28, count: 70},
			{from: 6, to: 21, count: 45},
			{from: 14, to: 22, count: 40},
			{from: 10, to: 26, count: 22},
		},
		total: 502,
	},
	14709037814934358540: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 21},
		},
		total: 21,
	},
	8345218066388553381: Book{
		moves: []shortMove{
			{from: 52, to: 45, count: 182},
		},
		total: 182,
	},
	14212068168608556130: Book{
		moves: []shortMove{
			{from: 10, to: 26, count: 148},
		},
		total: 148,
	},
	6634573336870380453: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 121},
		},
		total: 121,
	},
	8292814288991182844: Book{
		moves: []shortMove{
			{from: 6, to: 21, count: 56},
			{from: 13, to: 29, count: 27},
		},
		total: 83,
	},
	14264145469579210022: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 37},
		},
		total: 37,
	},
	16929927505613263231: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 37},
		},
		total: 37,
	},
	7530250141608134147: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 40},
		},
		total: 40,
	},
	7582145487114533971: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 22},
		},
		total: 22,
	},
	15457513895194024342: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 42},
		},
		total: 42,
	},
	17218845181714180084: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 28},
		},
		total: 28,
	},
	519880951014813005: Book{
		moves: []shortMove{
			{from: 2, to: 9, count: 858},
		},
		total: 858,
	},
	4559686210792690479: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 368},
			{from: 52, to: 44, count: 210},
			{from: 51, to: 35, count: 143},
			{from: 51, to: 43, count: 71},
			{from: 50, to: 34, count: 23},
		},
		total: 815,
	},
	3896481581505108126: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 178},
			{from: 52, to: 36, count: 153},
			{from: 62, to: 45, count: 101},
			{from: 50, to: 34, count: 69},
			{from: 54, to: 46, count: 40},
			{from: 52, to: 44, count: 33},
			{from: 51, to: 43, count: 21},
		},
		total: 595,
	},
	4271989409255176458: Book{
		moves: []shortMove{
			{from: 11, to: 19, count: 36},
			{from: 3, to: 10, count: 24},
			{from: 11, to: 27, count: 22},
		},
		total: 82,
	},
	11055715603404763611: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 21},
		},
		total: 21,
	},
	14281569926491073348: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 553},
			{from: 62, to: 45, count: 257},
			{from: 52, to: 36, count: 200},
			{from: 50, to: 34, count: 116},
			{from: 54, to: 46, count: 114},
			{from: 52, to: 44, count: 73},
			{from: 51, to: 43, count: 65},
			{from: 50, to: 42, count: 35},
			{from: 53, to: 37, count: 29},
			{from: 49, to: 41, count: 23},
		},
		total: 1465,
	},
	6687509113027777662: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 22},
		},
		total: 22,
	},
	17203337378853379824: Book{
		moves: []shortMove{
			{from: 14, to: 22, count: 29},
			{from: 6, to: 21, count: 21},
		},
		total: 50,
	},
	4359901050386863437: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 28},
		},
		total: 28,
	},
	1406874619383441684: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 26},
		},
		total: 26,
	},
	10517528860928392808: Book{
		moves: []shortMove{
			{from: 62, to: 45, count: 30},
			{from: 51, to: 43, count: 25},
			{from: 51, to: 35, count: 22},
		},
		total: 77,
	},
	11290513438659993596: Book{
		moves: []shortMove{
			{from: 1, to: 18, count: 67},
			{from: 6, to: 21, count: 23},
			{from: 1, to: 11, count: 20},
		},
		total: 110,
	},
	5135396313227408235: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 67},
		},
		total: 67,
	},
	18155078579001157095: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 32},
			{from: 2, to: 38, count: 26},
		},
		total: 58,
	},
	13803637700816989491: Book{
		moves: []shortMove{
			{from: 60, to: 62, count: 20},
		},
		total: 20,
	},
	14657078026839427792: Book{
		moves: []shortMove{
			{from: 12, to: 28, count: 70},
			{from: 6, to: 21, count: 57},
			{from: 14, to: 22, count: 45},
			{from: 1, to: 11, count: 24},
		},
		total: 196,
	},
	3445142071364069035: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 765},
			{from: 52, to: 36, count: 632},
			{from: 50, to: 34, count: 376},
			{from: 62, to: 45, count: 339},
			{from: 52, to: 44, count: 259},
			{from: 54, to: 46, count: 236},
			{from: 50, to: 42, count: 137},
			{from: 51, to: 43, count: 123},
			{from: 49, to: 41, count: 48},
			{from: 53, to: 37, count: 39},
		},
		total: 2954,
	},
	10098370877248517871: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 21},
		},
		total: 21,
	},
	14751760807717379281: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 72},
			{from: 10, to: 26, count: 23},
		},
		total: 95,
	},
	12441675744888340875: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 59},
			{from: 10, to: 26, count: 21},
		},
		total: 80,
	},
	13056780021991676305: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 126},
			{from: 10, to: 26, count: 48},
			{from: 9, to: 17, count: 44},
		},
		total: 218,
	},
	523339399611487007: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 132},
			{from: 6, to: 21, count: 28},
			{from: 9, to: 17, count: 22},
			{from: 10, to: 26, count: 21},
		},
		total: 203,
	},
	2492083325548414783: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 157},
			{from: 10, to: 26, count: 69},
			{from: 9, to: 17, count: 48},
			{from: 6, to: 21, count: 28},
		},
		total: 302,
	},
	14363958230478091807: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 102},
			{from: 52, to: 36, count: 73},
		},
		total: 175,
	},
	13126105616215140392: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 1538},
			{from: 62, to: 45, count: 622},
			{from: 54, to: 46, count: 408},
			{from: 52, to: 36, count: 399},
			{from: 50, to: 34, count: 264},
			{from: 50, to: 42, count: 218},
			{from: 52, to: 44, count: 171},
			{from: 51, to: 43, count: 144},
			{from: 55, to: 39, count: 85},
			{from: 53, to: 37, count: 68},
			{from: 57, to: 42, count: 46},
		},
		total: 3963,
	},
	1567902728827959404: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 28},
		},
		total: 28,
	},
	10930714570158059280: Book{
		moves: []shortMove{
			{from: 57, to: 42, count: 25},
		},
		total: 25,
	},
	6151382442137473618: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 198},
		},
		total: 198,
	},
	15505300296698889518: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 155},
		},
		total: 155,
	},
	3841403207901384456: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 129},
		},
		total: 129,
	},
	13203597635737409652: Book{
		moves: []shortMove{
			{from: 52, to: 36, count: 38},
			{from: 62, to: 45, count: 35},
			{from: 50, to: 42, count: 23},
		},
		total: 96,
	},
	13431608727669357024: Book{
		moves: []shortMove{
			{from: 6, to: 21, count: 21},
			{from: 10, to: 26, count: 20},
		},
		total: 41,
	},
	3230718888120473362: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 151},
		},
		total: 151,
	},
	12584795082554736750: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 116},
		},
		total: 116,
	},
	11816353135253152250: Book{
		moves: []shortMove{
			{from: 6, to: 21, count: 28},
		},
		total: 28,
	},
	133124376279902507: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 20},
		},
		total: 20,
	},
	11431709942539935132: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 361},
		},
		total: 361,
	},
	2077965215025466080: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 346},
		},
		total: 346,
	},
	3733811464214132409: Book{
		moves: []shortMove{
			{from: 6, to: 21, count: 92},
			{from: 11, to: 19, count: 71},
			{from: 10, to: 18, count: 44},
			{from: 10, to: 26, count: 41},
			{from: 1, to: 18, count: 41},
			{from: 11, to: 27, count: 25},
		},
		total: 314,
	},
	1264483158999591796: Book{
		moves: []shortMove{
			{from: 11, to: 19, count: 83},
			{from: 6, to: 21, count: 62},
			{from: 10, to: 26, count: 40},
			{from: 12, to: 28, count: 28},
			{from: 10, to: 18, count: 22},
		},
		total: 235,
	},
	4866689939720360063: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 25},
		},
		total: 25,
	},
	7822391173641301030: Book{
		moves: []shortMove{
			{from: 11, to: 27, count: 28},
		},
		total: 28,
	},
	12945796981840667557: Book{
		moves: []shortMove{
			{from: 61, to: 54, count: 105},
		},
		total: 105,
	},
	7692600509308976487: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 41},
		},
		total: 41,
	},
	13507273169351150012: Book{
		moves: []shortMove{
			{from: 5, to: 14, count: 577},
			{from: 6, to: 21, count: 25},
		},
		total: 602,
	},
	4107428715003502272: Book{
		moves: []shortMove{
			{from: 54, to: 46, count: 272},
			{from: 51, to: 35, count: 139},
			{from: 52, to: 44, count: 62},
			{from: 52, to: 36, count: 29},
			{from: 50, to: 34, count: 28},
			{from: 51, to: 43, count: 20},
		},
		total: 550,
	},
	17925040490297388515: Book{
		moves: []shortMove{
			{from: 51, to: 35, count: 36},
			{from: 52, to: 36, count: 35},
		},
		total: 71,
	},
}
