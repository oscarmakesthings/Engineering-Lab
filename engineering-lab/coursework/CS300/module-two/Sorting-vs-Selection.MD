Oscar Martinez
CS-300
Module Two: Selection Sort & Quicksort Algorithms

FIXME 1A:
FIXME 1A asked me to implement the selection sort logic over *bid.title*. The algorithm begins by declaring a 'min' variable where the smallest title will be stored. It continues by running an if statement to check for the size of *bids*. Beyond that we have a nested for loop which handles the main logic of the selection sort. The outer loop sets the positioning of the current element, making sure they are placed in the correct order. Everything 'placed' is considered sorted. The inner loop checks the rest of the unsorted portion to see if any title is smaller than the current one selected. If it is, *min* is updated and that element becomes the new placed minimum value. Then there is an if statement that verifies if something was found, and swaps it into current position if true.

FIXME 2A:
FIXME 2A asked me to implement the quick sort logic over _bid.title_. The idea is to pick a pivot title from the middle of the list, then move all bids with titles smaller than the pivot to the left and all bids with titles larger than the pivot to the right. This is done in the partition function, which uses two indexes (low and high) that move toward each other and swap elements when they are on the wrong side of the pivot. After that, quickSort uses the position returned by partition to split the list into two parts and then calls itself on the left and right sides. By repeating this process, the bids end up sorted by title.

FIXME 1B and 2B:
For FIXME 1B, I added code in menu option 3 to time the selectionSort function. I start the clock, call selectionSort(bids) to sort all the bids by title, then stop the clock and print how many ticks and seconds the sort took. For FIXME 2B, I did the same thing for quick sort in menu option 4. I start the clock, call quickSort(bids, 0, bids.size() - 1) if there are any bids, then stop the clock and print the total time so I can compare it with selection sort.