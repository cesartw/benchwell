namespace Benchwell.Utils {
	public bool fuzzy_match (string phrase, string item, out int score) {
		var terms = phrase.split ("|", -1);
		bool matched = false;
		score = 99999999;
		foreach (var term in terms) {
			int local_score = 0;
			if (fuzzy_term_match (term, item, out local_score)) {
				matched = true;
				if (local_score < score)
					score = local_score;
			}
		}
		if (!matched) {
			score = 0;
		}

		return matched;
	}

	private bool fuzzy_term_match (string term, string item, out int score) {
		score = 0;
		unichar termC;
		int termIndex = 0;
		int itemIndex = 0;
		int[] indices = {};
		if (term == item) {
			return true;
		}

		while (term.get_next_char (ref termIndex, out termC)) {
			unichar itemC;
			bool found = false;
			while (item.get_next_char (ref itemIndex, out itemC)) {
				if (itemC == termC) {
					found = true;
					indices += itemIndex;
					break;
				}
			}
			if (!found) {
				return false;
			}
		}

		for (var i = 1; i < indices.length; i++) {
			score += indices[i-1] + indices[i];
		}

		score += (term.length - item.length).abs ();

		return true;
	}
}
