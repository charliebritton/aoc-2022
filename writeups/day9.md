# Day 9
Once again a more difficult one than the previous days (making 7 and 8 seem really quite easy). This one was a very visually stimulating one and was easy enough to write. It helped to copy the output format verbaitem for verification and I can see myself automating this away with some sort of testing as it will show a diff where output is different.

Logging is getting to the point where a proper logger with a debug level or something could be useful as the actual test case takes a much longer time when outputting to console.

Would also be good to be able to specify the input filename at some point to make testing easier but once again there isn't too much point in doing all of this when I'm never going to look at this code again.

There's also a bug where the board is much taller than it needs to be (the output `#`s should be within `ROPE_LENGTH` `.`s of the edge).

Feeling good that I've managed to keep the streak but a bit down because all of the other people in ECS seem to be able to do it much faster than I can and are mostly finding it 'easy'... maybe it just means I'm using the wrong language?