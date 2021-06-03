from collections import Counter
import random


class HangManSolver:
	def __init__(self):
		with open("words.txt", "r") as f:
			self.words = f.readlines()
			for (i, word) in enumerate(self.words):
				self.words[i] = word.strip()
		random.shuffle(self.words)
		for (k, v) in enumerate(self.words):
			self.words[k] = v.lower()

		new = []
		for word in self.words:
			p = True
			for char in word:
				if ord(char) not in range(97, 122 + 1):
					p = False
			if p:
				new.append(word)
		self.words = new
		counter = Counter()
		total_letter_count = 0
		for word in self.words:
			letters = [char for char in word]
			total_letter_count += len(word)
			counter.update(letters)

		self.letter_probs = [(i, counter[i] / total_letter_count * 100.0) for i in counter]

	def find_matches(self, correct, incorrect):
		min_score = sum([1 for i in correct if i != "_"])
		matches = []
		for word in self.words:
			if len(word) == len(correct):
				score = 0
				for i in range(len(correct)):
					if word[i] == correct[i]:
						score += 1
				if score >= min_score:
					p = True
					for i in incorrect:
						if i in word:
							p = False
					if p:
						matches.append(word)
		return matches

	def find_probabilities(self, correct, list):
		missing = 0
		for (k, v) in enumerate(correct):
			if v == "_":
				missing = k
				break
		letters = []
		for i in list:
			letters.append(i[missing])

		c = Counter(letters)
		probs = [(i, c[i] / len(letters) * 100.0) for i in c]

		for i in range(len(probs)):
			for j in range(len(self.letter_probs)):
				if probs[i][0] == self.letter_probs[j][0]:
					s = probs[i][1] + self.letter_probs[j][1]
					s /= 2
					probs[i] = (probs[i][0], s)

		probs.sort(key=lambda x: x[1])
		probs.reverse()

		pruned_probs = []
		for (l, p) in probs:
			if l not in correct:
				pruned_probs.append((l, p))
		return pruned_probs

"""
solver = HangManSolver()

max_steps = 0
max_word = 0
n_trials = 50
failed = 0
max_failed = 0
game_failed_steps = 8
print(len(solver.words))
for word in solver.words[0:n_trials]:
	correct = ""
	incorrect = ""
	steps = 0
	score = 0

	for i in word:
		correct += "_"
	while "_" in correct:
		matches = solver.find_matches(correct, incorrect)

		guess = solver.find_probabilities(correct, matches)[0][0]

		if len(matches) == 1 and matches[0] == word:
			correct = word

		if guess in word:
			for (k, v) in enumerate(word):
				if guess == v:
					l = list(correct)
					l[k] = guess
					correct = "".join(l)
		else:
			incorrect += guess
			score += 1

		print(word, correct, guess, score)
		steps += 1
	if steps > max_steps:
		max_steps = steps
		max_word = word
	if score >= game_failed_steps:
		failed += 1
	if score > max_failed:
		max_failed = score

print("most steps:", max_steps)
print("most word:", max_word)
print("most failed:", max_failed)
print("percent failed:", failed / n_trials * 100)
"""
#"""
solver = HangManSolver()
correct = "pop"
incorrect = "ascetibdnlr"

words = solver.find_matches(correct, incorrect)
#print("possible words:", words)
#print()
guess = solver.find_probabilities(correct, words)[0][0]
print("possible next letters:", guess)
#"""

