# time-parser
Basic time parser using a Go PEG that returns minutes after midnight

This PEG was taken downloaded from https://github.com/pointlander/peg

## Results
I wrote grammar rules in BNF that defined basic standard and military time formats such as:

` 4pm, 7:38pm, 23:42, 3:16, 3:16am `

I took these rules and wrote them in their PEG forms as .peg file, which the generates a parser that I'm able to pass any string.

The expected out for a valid time like the above is the number of minutes after midnight. 


This exercise was taken from The Pragmatic Programmer and was my first time working with a generator like this and first time in years looking at BNF.
It really made me appreciate and consider more thoughtfully the work that goes into writing a programming language and the amount of time that is likely saved by using a PEG

## Future 

If I get the chance I would like to improve some of the logic being using in the grammar file. I would also like to spend time thinking of the actual real world applications of a PEG 
