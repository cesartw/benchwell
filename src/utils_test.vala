class TestUtils : GLib.TestSuite {
	public TestUtils () {
		// assign a name for this class
		base ("TestExample");
		// add test methods
		add (new TestFuzzy ());
	}

}

class TestFuzzy : GLib.TestCase {
	public TestFuzzy () {
		base ("fuzzy", set_up, tear_down, test_samples);
	}

	public void set_up () {
	}

	public void test_samples () {
		var terms = "abcdfg";

		string[] search_terms = {
			"adg",
			"bcd",
			"acf",
			"ach",
			"cda"
		};

		int[] scores = {
			2+1+3,
			0+3,
			1+1+3,
			0,
			0
		};

		bool[] founds = {
			true,
			true,
			true,
			false,
			false
		};

		for (var i = 0; i < search_terms.length; i++) {
			int score = 0;
			var found = Benchwell.Utils.fuzzy_match (search_terms[i], terms, out score);

			assert_cmpint (score, GLib.CompareOperator.EQ, scores[i]);
			assert (founds[i] == found);
		}
	}

	public void tear_down () {
	}
}
