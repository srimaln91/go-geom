package geos

import (
	"regexp"
	"testing"
)

var WKT = "LINESTRING(79.856178 6.911853, 79.856382 6.911449, 79.856645 6.910948, 79.85709 6.910065, 79.857331 6.909487, 79.857387 6.909367, 79.857499 6.909128, 79.857976 6.908112, 79.858588 6.906801, 79.859 6.906006, 79.859183 6.905839, 79.859722 6.905565, 79.859839 6.905507, 79.860013 6.905459, 79.860314 6.905412, 79.860492 6.905392, 79.86066 6.905374, 79.860996 6.905357, 79.86121 6.905347, 79.861278 6.905347, 79.861595 6.905353, 79.861619 6.905352, 79.861875 6.905382, 79.862101 6.90543, 79.862288 6.905479, 79.862517 6.90539, 79.862581 6.905314, 79.862657 6.905189, 79.862679 6.905029, 79.862653 6.904882, 79.862586 6.904199, 79.862554 6.903853, 79.862539 6.903703, 79.862527 6.903659, 79.862495 6.903508, 79.862386 6.903139, 79.862051 6.902145, 79.861899 6.901833, 79.860774 6.899256, 79.860674 6.89877, 79.860163 6.89772, 79.860179 6.897692, 79.860189 6.89766, 79.860192 6.897627, 79.860188 6.897595, 79.860178 6.897564, 79.860161 6.897535, 79.860138 6.897511, 79.860111 6.897492, 79.859987 6.897294, 79.859992 6.89728, 79.859995 6.897265, 79.859997 6.89725, 79.859997 6.897235, 79.859995 6.89722, 79.859992 6.897205, 79.859988 6.89719, 79.859973 6.897161, 79.859951 6.897135, 79.859925 6.897017, 79.859913 6.896801, 79.859948 6.896713, 79.859992 6.896627, 79.860084 6.896499, 79.860162 6.896395, 79.860241 6.896296, 79.860332 6.896182, 79.860404 6.896078, 79.860468 6.895988, 79.860541 6.895861, 79.860742 6.895464, 79.860809 6.895337, 79.860891 6.895185, 79.86104 6.894948, 79.861146 6.894817, 79.861288 6.894679, 79.861463 6.894502, 79.86153 6.894417, 79.861613 6.894327, 79.861684 6.894205, 79.861725 6.894083, 79.861773 6.893928, 79.861902 6.893542, 79.862071 6.892919, 79.862108 6.89277, 79.862154 6.89264, 79.862213 6.892506, 79.862299 6.892317, 79.862417 6.892244, 79.862491 6.892137, 79.862578 6.892026, 79.862596 6.892003, 79.862692 6.89188, 79.862818 6.891737, 79.862959 6.891598, 79.863101 6.891475, 79.863271 6.891329, 79.863382 6.891226, 79.863508 6.891107, 79.86355 6.89106, 79.863589 6.89101, 79.863627 6.890952, 79.863662 6.890887, 79.863693 6.890822, 79.863721 6.890756, 79.863749 6.890677, 79.863776 6.890582, 79.863802 6.890448, 79.863822 6.890353, 79.863831 6.890314, 79.863864 6.890182, 79.863956 6.889815, 79.864022 6.889543, 79.864115 6.88916, 79.864239 6.888561, 79.864302 6.888267, 79.864269 6.88813, 79.864377 6.88784, 79.864531 6.887544, 79.864612 6.887389, 79.864718 6.887183, 79.86485 6.886945, 79.865007 6.886708, 79.865147 6.886517, 79.86521 6.886434, 79.86525 6.886375, 79.865612 6.886062, 79.865788 6.885912, 79.865981 6.885763, 79.866178 6.885644, 79.866411 6.885508, 79.866649 6.885361, 79.866932 6.885135, 79.867028 6.88507, 79.867311 6.884844, 79.867529 6.884661, 79.867656 6.884521, 79.867838 6.884294, 79.868023 6.884063, 79.868219 6.883822, 79.868302 6.883727, 79.868478 6.883514, 79.868629 6.883351, 79.868986 6.882925, 79.86912 6.882774, 79.869247 6.882629, 79.869352 6.882498, 79.86942 6.8824, 79.869474 6.88232, 79.869554 6.882198, 79.869836 6.881867, 79.86986 6.881839, 79.869876 6.881805, 79.869882 6.881769, 79.869973 6.881601, 79.870037 6.881464, 79.870062 6.881386, 79.870112 6.881297, 79.87128 6.879041, 79.871435 6.878831, 79.871536 6.878742, 79.87163 6.878681, 79.871723 6.878629, 79.871845 6.878589, 79.871967 6.878556, 79.872073 6.878528, 79.872569 6.87856, 79.872715 6.878585, 79.873133 6.878678, 79.873533 6.878738, 79.873625 6.87874, 79.873707 6.878745, 79.873814 6.878748, 79.873918 6.878743, 79.874105 6.87874, 79.87454 6.878695, 79.875431 6.878559, 79.8764 6.87834, 79.876511 6.878315, 79.876854 6.878199, 79.877169 6.878105, 79.877283 6.878079, 79.877585 6.878013, 79.878168 6.877869, 79.878302 6.877841, 79.879021 6.877657, 79.879192 6.877605, 79.879406 6.877519, 79.879486 6.877484, 79.879574 6.877439, 79.879803 6.877278, 79.880019 6.877099, 79.880203 6.876916, 79.880355 6.876765, 79.88043 6.876673, 79.880564 6.876498, 79.880741 6.87621, 79.880893 6.875897, 79.881106 6.875408, 79.881294 6.875007, 79.881591 6.874373, 79.881769 6.873994, 79.882007 6.873543, 79.882081 6.873442, 79.88229 6.873195, 79.882422 6.873077, 79.882548 6.872979, 79.883188 6.872585, 79.883421 6.872454, 79.883603 6.87235, 79.884041 6.872111, 79.885476 6.871327, 79.886506 6.870761, 79.886619 6.870699, 79.886736 6.870635, 79.886749 6.870628, 79.887002 6.870498, 79.887085 6.870452, 79.88715 6.870438, 79.887257 6.870405, 79.887406 6.870339, 79.887501 6.870295, 79.887872 6.870135, 79.888305 6.869959, 79.888594 6.86984, 79.889042 6.869647, 79.889479 6.869472, 79.889656 6.869398, 79.889842 6.869302, 79.889999 6.8692, 79.890292 6.869087, 79.890554 6.869008, 79.890923 6.868875, 79.891457 6.868663, 79.891688 6.868575, 79.891809 6.868509, 79.891887 6.868467, 79.891973 6.868413, 79.892043 6.868369, 79.892106 6.868323, 79.892185 6.868266, 79.892264 6.86821, 79.892438 6.868076, 79.892608 6.867931, 79.892774 6.867781, 79.892871 6.86767, 79.893024 6.867449, 79.893139 6.867252, 79.893281 6.867009, 79.893435 6.866701, 79.893697 6.866198, 79.893779 6.866053, 79.893937 6.865895, 79.894026 6.86582, 79.894132 6.865749, 79.894268 6.865695, 79.894396 6.865659, 79.894562 6.865615, 79.894659 6.865593, 79.894775 6.865566, 79.895049 6.865533, 79.895404 6.865502, 79.895554 6.865485, 79.895732 6.865466, 79.896001 6.865399, 79.896322 6.865275, 79.896616 6.865174, 79.897396 6.864723, 79.897901 6.864443, 79.898059 6.864353, 79.898312 6.864219, 79.898802 6.863928, 79.899398 6.863662, 79.899562 6.863586, 79.899716 6.863523, 79.899869 6.863466, 79.899979 6.86341, 79.900444 6.86329, 79.901439 6.862979, 79.901607 6.862953, 79.901705 6.862916, 79.903737 6.861928, 79.903888 6.861862, 79.905027 6.861365, 79.90518 6.861279, 79.905228 6.861252, 79.905494 6.861049, 79.905736 6.860815, 79.90594 6.860563, 79.906286 6.860145, 79.906334 6.860089, 79.907342 6.859069, 79.907619 6.858788, 79.907847 6.858545, 79.908087 6.858291, 79.908303 6.85806, 79.90843 6.857899, 79.908544 6.857777, 79.908808 6.857512, 79.908916 6.857415, 79.909046 6.857315, 79.909117 6.857264, 79.909219 6.857197, 79.909311 6.857136, 79.909434 6.857056, 79.909699 6.856899, 79.910407 6.856559, 79.910538 6.8565, 79.911301 6.856143, 79.911565 6.85602, 79.911886 6.855863, 79.912374 6.85561, 79.913431 6.855063, 79.91394 6.854744, 79.914096 6.854624, 79.914262 6.854494, 79.914349 6.854403, 79.914482 6.854262, 79.91502 6.853648, 79.915086 6.853568, 79.915368 6.853236, 79.915687 6.852888, 79.915897 6.852713, 79.915999 6.852641, 79.916067 6.852596, 79.916167 6.852548, 79.91627 6.852508, 79.916442 6.852453, 79.916572 6.85243, 79.91664 6.85242, 79.916733 6.852423, 79.917081 6.852426, 79.917451 6.852428, 79.917624 6.852421, 79.917734 6.852413, 79.917856 6.852405, 79.918102 6.852382, 79.918417 6.852324, 79.918619 6.852275, 79.918865 6.852208, 79.91904 6.852147, 79.91944 6.852002, 79.919737 6.851894, 79.919839 6.851856, 79.920607 6.851535, 79.920662 6.851511, 79.920787 6.851444, 79.920921 6.851364, 79.921043 6.851286, 79.921162 6.85121, 79.921265 6.851137, 79.922242 6.850478, 79.922767 6.85015, 79.923306 6.849811, 79.923574 6.849629, 79.923974 6.849351, 79.924288 6.849133, 79.924714 6.848856, 79.924897 6.848641, 79.924954 6.848576, 79.92512 6.848419, 79.92522 6.848352, 79.925419 6.848205, 79.925541 6.848092, 79.925588 6.848073, 79.925888 6.847744, 79.926268 6.847406, 79.926597 6.847207, 79.927638 6.846543, 79.927658 6.846529, 79.927688 6.846511, 79.928752 6.845895, 79.929484 6.845431, 79.929914 6.845136, 79.93027 6.845033, 79.930424 6.844963, 79.930612 6.844904, 79.930852 6.844861, 79.931028 6.844825, 79.931366 6.844796, 79.931734 6.844777, 79.931959 6.844774, 79.932085 6.844764, 79.932174 6.844759, 79.932524 6.844725, 79.932995 6.844681, 79.933278 6.844649, 79.93351 6.844646, 79.933748 6.84465, 79.933876 6.844656, 79.934235 6.8447, 79.934379 6.844747, 79.934613 6.84485, 79.935093 6.845029, 79.935258 6.845099, 79.93547 6.845173, 79.935708 6.845272, 79.93594 6.845362, 79.936151 6.845431, 79.936355 6.845474, 79.936511 6.845488, 79.936608 6.845498, 79.937052 6.845548, 79.937488 6.845584, 79.937686 6.845574, 79.937832 6.845562, 79.938144 6.845525, 79.938432 6.84549, 79.938735 6.845437, 79.938982 6.845389, 79.939274 6.845327, 79.940122 6.845123, 79.940312 6.845086, 79.940705 6.845022, 79.94081 6.845005, 79.94089 6.844993, 79.941036 6.844974, 79.941226 6.844955, 79.941445 6.844964, 79.94168 6.844994, 79.942037 6.845043, 79.942411 6.845097, 79.942856 6.845172, 79.943309 6.845266, 79.943809 6.845378, 79.944537 6.845571, 79.94554 6.845821, 79.946188 6.845987, 79.946731 6.84613, 79.947004 6.84621, 79.947541 6.846356, 79.947809 6.846411, 79.947944 6.846431, 79.948081 6.846438, 79.948253 6.846439, 79.948395 6.846431, 79.948501 6.846399, 79.948602 6.846374, 79.948772 6.846331, 79.948852 6.846312, 79.949425 6.84621, 79.949725 6.846167, 79.950033 6.846134, 79.951147 6.846018, 79.95156 6.845958, 79.952068 6.845847, 79.952413 6.845726, 79.952832 6.845568, 79.953245 6.845397, 79.954034 6.845013, 79.954656 6.844723, 79.95475 6.844684, 79.955182 6.844496, 79.956137 6.844081, 79.956603 6.843851, 79.95685 6.843754, 79.957134 6.843669, 79.957322 6.843633, 79.957501 6.843618, 79.957783 6.843601, 79.958584 6.843601, 79.959064 6.843579, 79.959558 6.843532, 79.959802 6.843513, 79.959945 6.843495, 79.960142 6.843445, 79.960402 6.843362, 79.960731 6.843244, 79.961022 6.843093, 79.961521 6.842774, 79.961753 6.842637, 79.962021 6.842458, 79.962172 6.842353, 79.96227 6.842286, 79.962451 6.84217, 79.962638 6.842047, 79.962794 6.841957, 79.962913 6.841898, 79.963077 6.841832, 79.963326 6.841762, 79.963455 6.841743, 79.963575 6.841727, 79.963777 6.841705, 79.963977 6.84169, 79.964166 6.841748, 79.964237 6.841783, 79.964279 6.841804, 79.964383 6.841865, 79.964417 6.841901, 79.964458 6.841923, 79.964507 6.841933, 79.964556 6.841927, 79.964601 6.841906, 79.964637 6.841871, 79.964792 6.841743, 79.964846 6.841691, 79.965035 6.841559, 79.965197 6.841445, 79.965363 6.841338, 79.965496 6.841249, 79.965684 6.841172, 79.965933 6.841052, 79.966153 6.840945, 79.966349 6.840847, 79.966873 6.840598, 79.967273 6.84038, 79.967406 6.840301, 79.967555 6.840215, 79.96773 6.840112, 79.967995 6.839939, 79.968534 6.839609, 79.968682 6.839543, 79.96876 6.839514, 79.968848 6.839492, 79.96919 6.839417, 79.969705 6.839366, 79.970086 6.839356, 79.970558 6.839394, 79.970905 6.839424, 79.971205 6.839432, 79.971263 6.839427, 79.971504 6.839409, 79.97177 6.839355, 79.971963 6.839317, 79.972339 6.839214, 79.972747 6.839124, 79.973069 6.839079, 79.973267 6.839053, 79.973448 6.83906, 79.973684 6.839065, 79.973803 6.839081, 79.97392 6.839103, 79.974327 6.839206, 79.974713 6.839313, 79.974933 6.839386, 79.975056 6.839428, 79.975297 6.839484, 79.975553 6.839519, 79.975832 6.839528, 79.976059 6.839523, 79.976453 6.839501, 79.976905 6.839457, 79.977084 6.839426, 79.977217 6.839394, 79.977346 6.839363, 79.977645 6.839248, 79.977969 6.8391, 79.978276 6.838948, 79.978446 6.838882, 79.978608 6.838832, 79.97875 6.838796, 79.978902 6.838771, 79.979038 6.838745, 79.979302 6.838704, 79.979445 6.838688, 79.97959 6.838664, 79.979719 6.838658, 79.980024 6.838649, 79.980099 6.838644, 79.980205 6.83863, 79.980458 6.838592, 79.980777 6.838543, 79.980814 6.838537, 79.980835 6.838525, 79.98085 6.838502, 79.980923 6.838246, 79.980941 6.838184, 79.980944 6.838143, 79.980937 6.838117, 79.980914 6.838106, 79.98088 6.838104, 79.980829 6.838118, 79.980724 6.838159, 79.980541 6.838238, 79.980406 6.838282, 79.980216 6.838349, 79.980212 6.838443, 79.980091 6.838486, 79.97987 6.838569, 79.979838 6.83858)"

// Clear a set of Geoms from memory
func cleanup(geoms ...*Geom) {
	for _, geom := range geoms {
		geom.Destroy()
	}
}

func TestFromWKT(t *testing.T) {

	geom := FromWKT(WKT)

	if geom.cGeom == nil {
		t.Errorf("Error: CreateFromWKT error")
	}
}

func TestBuffer(t *testing.T) {

	geom := FromWKT("POINT (0 0)")
	geom.Buffer(1)

	if geom.cGeom == nil {
		t.Errorf("Error: Buffer() error")
	}

}

func TestToWKT(t *testing.T) {

	wkt := "POINT (0 0)"
	expectedWKT := "POINT (0.0000000000000000 0.0000000000000000)"

	geom := FromWKT(wkt)

	if geom.ToWKT() != expectedWKT {
		t.Errorf("Error: ToWKT(%s) error", wkt)
	}
}

func TestSRID(t *testing.T) {
	geom := FromWKT(WKT)
	geom.SetSRID(4326)
	srid := geom.GetSRID()

	if srid != 4326 {
		t.Errorf("Error: SRID(%s) error", WKT)
	}
}

func BenchmarkBuffer(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			geom := FromWKT(WKT)
			geom.Buffer(1)
		}

	})

}

func TestSimplify(t *testing.T) {

	geom := FromWKT(WKT)

	geom.Simplify(0.01)

	resultWKT := geom.ToWKT()
	expectedWKT := "LINESTRING (79.8561779999999999 6.9118529999999998, 79.8716299999999961 6.8786810000000003, 79.9299139999999966 6.8451360000000001, 79.9798380000000009 6.8385800000000003)"

	if resultWKT != expectedWKT {
		t.Errorf("Error: Simplify(%s) error", resultWKT)
	}
}

func TestSimplifyPreserveTopology(t *testing.T) {

	geom := FromWKT(WKT)

	geom.SimplifyPreserveTopology(0.01)

	resultWKT := geom.ToWKT()
	expectedWKT := "LINESTRING (79.8561779999999999 6.9118529999999998, 79.8716299999999961 6.8786810000000003, 79.9299139999999966 6.8451360000000001, 79.9798380000000009 6.8385800000000003)"

	if resultWKT != expectedWKT {
		t.Errorf("Error: SimplifyPreserveTopology(%s) error", resultWKT)
	}
}

func TestReverse(t *testing.T) {
	wkt := "LINESTRING (30 10, 10 30, 40 40)"
	geom := FromWKT(wkt)

	geom.Reverse()

	resultWKT := geom.ToWKT()

	if resultWKT != "LINESTRING (40.0000000000000000 40.0000000000000000, 10.0000000000000000 30.0000000000000000, 30.0000000000000000 10.0000000000000000)" {
		t.Errorf("Error: ToWKT(%s) error", resultWKT)
	}
}

func TestVersion(t *testing.T) {
	version := Version()

	matched, err := regexp.MatchString(`^3\.\d+\.\d+-CAPI-\d+\.\d+\.\d+.+$`, version)

	if err != nil {
		t.Errorf("Error: Version fetch error, %s", err)
	}

	if !matched {
		t.Errorf("Error: Version %s is invalid", version)
	}
}

func TestUnion(t *testing.T) {

	geom1 := FromWKT("LINESTRING(79.856178 6.911853, 79.856382 6.911449, 79.856645 6.910948)")
	geom2 := FromWKT("LINESTRING(79.856178 6.911853, 79.856382 6.911449, 79.856645 6.910948)")

	union := geom1.Union(geom2)

	resultWKT := union.ToWKT()
	expectedWKT := "MULTILINESTRING ((79.8561779999999999 6.9118529999999998, 79.8563819999999964 6.9114490000000002), (79.8563819999999964 6.9114490000000002, 79.8566450000000003 6.9109480000000003))"

	if resultWKT != expectedWKT {
		t.Errorf("Error: Invalid Buffer %s", resultWKT)
	}

	//Cleanup
	cleanup(geom1, geom2, union)
}

func TestIntersection(t *testing.T) {

	geom1 := FromWKT("POLYGON((79.856178 6.911853,79.85598771527475 6.911267935124347,79.85636717179295 6.911461947090437,79.856178 6.911853))")
	geom2 := FromWKT("POLYGON((79.85599371221633 6.911822745245366,79.85623376992316 6.911505881917993,79.85612245824905 6.911199669256946,79.85599371221633 6.911822745245366))")

	intersection := geom1.Intersection(geom2)

	resultWKT := intersection.ToWKT()
	expectedWKT := "POLYGON ((79.8560614850391488 6.9114947536352744, 79.8561157638370389 6.9116616436370890, 79.8562337699231648 6.9115058819179929, 79.8561836990276532 6.9113681394729278, 79.8560968252871675 6.9113237218799171, 79.8560614850391488 6.9114947536352744))"

	if resultWKT != expectedWKT {
		t.Errorf("Error: Invalid Intersection %s", resultWKT)
	}

	//Cleanup
	cleanup(geom1, geom2, intersection)
}

func TestIntersects(t *testing.T) {
	geom1 := FromWKT("POLYGON((79.856178 6.911853,79.85598771527475 6.911267935124347,79.85636717179295 6.911461947090437,79.856178 6.911853))")
	geom2 := FromWKT("POLYGON((79.85599371221633 6.911822745245366,79.85623376992316 6.911505881917993,79.85612245824905 6.911199669256946,79.85599371221633 6.911822745245366))")

	intersects, _ := geom1.Intersects(geom2)

	if intersects == false {
		t.Errorf("Error: Intersects")
	}

	//Cleanup
	cleanup(geom1, geom2)

}

func TestDisjoint(t *testing.T) {

	geom1 := FromWKT("POINT(0 0)")
	geom2 := FromWKT("LINESTRING ( 2 0, 0 2 )")
	geom3 := FromWKT("LINESTRING ( 0 0, 0 2 )")

	disjoint, _ := geom1.Disjoints(geom2)

	if disjoint == false {
		t.Errorf("Error: Disjoint")
	}

	disjoint, _ = geom1.Disjoints(geom3)

	if disjoint == true {
		t.Errorf("Error: Disjoint")
	}

	//Cleanup
	cleanup(geom1, geom2, geom3)
}

func TestTouches(t *testing.T) {
	geom1 := FromWKT("POINT(1 1)")
	geom2 := FromWKT("POINT(0 2)")
	geom3 := FromWKT("LINESTRING(0 0, 1 1, 0 2)")
	geom4 := FromWKT("LINESTRING(0 0, 1 1, 0 2)")

	touches, _ := geom1.Touches(geom3)

	if touches == true {
		t.Errorf("Error: Touches")
	}

	touches, _ = geom2.Touches(geom4)

	if touches == false {
		t.Errorf("Error: Touches")
	}

	//Cleanup
	cleanup(geom1, geom2, geom3, geom4)
}

func TestWithin(t *testing.T) {

	geom1 := FromWKT("POLYGON((20 15,23 11,26 16,20 15))")
	geom2 := FromWKT("POLYGON((8 52,21 50,11 46,8 52))")
	geom3 := FromWKT("POLYGON((10 10,25.703125 7,40 10,33 19,24 26,15 20,10 10))")

	disjoint, _ := geom1.Within(geom3)

	if disjoint == false {
		t.Errorf("Error: Within")
	}

	disjoint, _ = geom2.Within(geom3)

	if disjoint == true {
		t.Errorf("Error: Within")
	}

	//Cleanup
	cleanup(geom1, geom2, geom3)
}
