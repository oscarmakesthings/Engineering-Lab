## Overview

`VectorSorting.cpp` is a small console program that:
- Loads auction bids from a CSV file into a `vector<Bid>`.
- Displays bids.
- Sorts bids by `title` using two algorithms you must implement: **selection sort** and **quick sort**.
- Measures and prints how long each sort takes.

## Main data structures and functions

- **`struct Bid`**: holds a single bid:
  - `bidId`, `title`, `fund` (all `string`), and `amount` (`double`).
- **`displayBid(const Bid&)`**: prints one bid on a single line.
- **`getBid()`**: reads a bid from standard input (not used in the menu, but useful as a helper).
- **`loadBids(string csvPath)`**:
  - Uses `csv::Parser` to read each row from the CSV file.
  - Builds a `Bid` from columns (title, id, fund, amount) and pushes it into `vector<Bid>`.
  - Returns the full `vector<Bid>`.
- **`strToDouble(string, char)`**:
  - Removes a given character (e.g. `$`) and converts the remaining string to `double`.

## Control flow in `main`

- Reads an optional CSV path from `argv` (default is `eBid_Monthly_Sales.csv`).
- Maintains a `vector<Bid> bids` and a `clock_t ticks` for timing.
- Presents a menu in a loop:
  - **1. Load Bids**: calls `loadBids`, prints how many bids were read, and prints elapsed time.
  - **2. Display All Bids**: loops through `bids` and calls `displayBid` on each.
  - **3. Selection Sort All Bids**: **(you must implement – see FIXME 1a/1b)**.
  - **4. Quick Sort All Bids**: **(you must implement – see FIXME 2a/2b)**.
  - **9. Exit**: ends the loop and prints “Good bye.”.

## Sorting responsibilities (FIXMEs)

There are four main FIXMEs:

### FIXME (1a): Implement `selectionSort(vector<Bid>& bids)`

Goal: sort the `bids` vector **in place** by `bid.title` using **selection sort**.

What you need to do:
- Declare `size_t size = bids.size();`.
- Outer loop: `for (size_t pos = 0; pos < size - 1; ++pos)`:
  - Set `min = pos;`.
  - Inner loop: `for (size_t j = pos + 1; j < size; ++j)`:
    - If `bids[j].title < bids[min].title`, update `min = j;`.
  - After the inner loop, swap the smallest element found with the one at `pos`
    (e.g. `swap(bids[pos], bids[min]);`).

### FIXME (1b): Call `selectionSort` from the menu and time it

In the `switch (choice)` inside `main`, for **case 3** you should:
- Start the timer: `ticks = clock();`.
- Call `selectionSort(bids);`.
- Stop the timer: `ticks = clock() - ticks;`.
- Print:
  - number of bids (`bids.size()`),
  - clock ticks,
  - seconds (`ticks * 1.0 / CLOCKS_PER_SEC`).
- Optionally display the first few bids to confirm the sort worked.

### FIXME (2a): Implement quick sort (`partition` and `quickSort`)

You must implement:
- `int partition(vector<Bid>& bids, int begin, int end)`:
  - Initialize `int low = begin; int high = end;`.
  - Compute middle: `int mid = begin + (end - begin) / 2;`.
  - Set `string pivot = bids[mid].title;`.
  - Loop until `done`:
    - Increment `low` while `bids[low].title < pivot`.
    - Decrement `high` while `pivot < bids[high].title`.
    - If `low >= high`, return `high` (partition index).
    - Otherwise, `swap(bids[low], bids[high]); ++low; --high;`.
- `void quickSort(vector<Bid>& bids, int begin, int end)`:
  - Base case: if `begin >= end`, return (0 or 1 element).
  - Call `int mid = partition(bids, begin, end);`.
  - Recursively sort the two partitions:
    - `quickSort(bids, begin, mid);`
    - `quickSort(bids, mid + 1, end);`

This will sort the vector by `title` using quick sort, with average complexity \(O(n \log n)\).

### FIXME (2b): Call `quickSort` from the menu and time it

In the `switch (choice)` inside `main`, for **case 4** you should:
- Start the timer: `ticks = clock();`.
- Call `quickSort(bids, 0, bids.size() - 1);` (only if `bids` is not empty).
- Stop the timer and print:
  - number of bids,
  - clock ticks,
  - seconds (as with selection sort).
- Optionally display the first few bids to see that they’re sorted.

## Summary

- The program loads bid data into a `vector<Bid>`, then lets you sort that vector by `title` using two different algorithms.
- Your main tasks are to:
  - Implement **selection sort** and **quick sort** over `bids` (by `title`).
  - Wire them into the menu options **3** and **4** with timing output.
- Once done, you can compare the timing results for selection sort vs quick sort on the same data set.

