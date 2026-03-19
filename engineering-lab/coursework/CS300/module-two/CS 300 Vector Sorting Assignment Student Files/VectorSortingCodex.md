# VectorSortingCodex Notes

## Zybooks alignment

I updated `VectorSortingCodex.cpp` to follow the standard approach typically taught in Zybooks for `selectionSort` and `quickSort`:

- Selection sort scans the unsorted portion of the vector to find the minimum element.
- After finding the minimum, it swaps that element into the current position.
- Quicksort uses a helper `partition()` function.
- The partition function chooses the middle element as the pivot.
- Two indices, `low` and `high`, move toward each other until elements need to be swapped.
- Quicksort then recursively sorts the left and right partitions.

This matches the common Zybooks pattern for these algorithms in C++.

## What I changed in the code

I created `VectorSortingCodex.cpp` as a duplicate of the original assignment file and implemented only the missing `FIXME` sections there.

### `FIXME (1a)` Selection sort

I completed `selectionSort(vector<Bid>& bids)` so it now:

- Stores the vector size in `size`.
- Returns immediately if the vector has 0 or 1 elements.
- Uses `pos` to mark the boundary between the sorted and unsorted parts of the vector.
- Uses `min` to track the index of the smallest `title` in the unsorted portion.
- Swaps the smallest found element into the current position.

This sorts bids alphabetically by `title`.

### `FIXME (2a)` Quick sort

I completed:

- `partition(vector<Bid>& bids, int begin, int end)`
- `quickSort(vector<Bid>& bids, int begin, int end)`

The quicksort implementation now works like this:

- `partition()` sets `low = begin` and `high = end`.
- It picks the middle element's `title` as the pivot.
- It increments `low` while the current title is less than the pivot.
- It decrements `high` while the pivot is less than the current title.
- If `low` is still left of `high`, it swaps those elements and keeps moving inward.
- When `low >= high`, it returns `high` as the partition point.
- `quickSort()` stops when `begin >= end`.
- Otherwise, it partitions the vector and recursively sorts both halves.

### `FIXME (1b)` and `FIXME (2b)` Menu options

I added the missing menu actions in `main()`:

- `case 3` runs selection sort and prints elapsed clock ticks and seconds.
- `case 4` runs quicksort and prints elapsed clock ticks and seconds.

I also kept a guard so quicksort only runs when the vector is not empty.

## Selection sort vs. quicksort

### Selection sort

Selection sort repeatedly finds the smallest element in the unsorted part of the vector and places it into the next sorted position.

Characteristics:

- Easy to understand and implement.
- Uses nested loops.
- Performs well only on small datasets.
- Time complexity is `O(n^2)` in average and worst cases.
- It does not benefit much from partially sorted data.

In this assignment, selection sort is useful because it clearly shows the basic idea of comparing titles and moving the smallest one forward.

### Quicksort

Quicksort chooses a pivot, partitions the vector around that pivot, and then recursively sorts the two resulting parts.

Characteristics:

- More efficient than selection sort for large datasets.
- Uses divide-and-conquer recursion.
- Average time complexity is `O(n log n)`.
- Worst-case time complexity is `O(n^2)`.
- In practice, it is usually much faster than selection sort on large input files.

In this assignment, quicksort is better suited for sorting a large list of bids because it reduces the amount of total work compared with checking every remaining element on every pass.

## Summary comparison

- Selection sort is simpler and easier to trace by hand.
- Quicksort is more efficient for larger collections.
- Selection sort always does roughly the same amount of work.
- Quicksort usually does much less work on average.
- Selection sort is good for learning the basics of sorting.
- Quicksort is better when performance matters.

## Result

`VectorSortingCodex.cpp` now:

- Loads bids from the CSV file.
- Displays all bids.
- Sorts bids by `title` using selection sort.
- Sorts bids by `title` using quicksort.
- Reports timing for both algorithms.
