//============================================================================
// FIXME Implementations for VectorSorting.cpp
// CS 300 Vector Sorting Assignment
//
// Copy each section into the corresponding location in VectorSorting.cpp.
// This file is for reference only and is not compiled on its own.
//============================================================================

//----------------------------------------------------------------------------
// FIXME (2a) — Quick sort: partition and quickSort
// Replace the stub bodies in VectorSorting.cpp (lines ~123–165).
//----------------------------------------------------------------------------

/**
 * Partition the vector of bids into two parts, low and high.
 * Uses the middle element as pivot; elements with title < pivot end up in the
 * low part, elements with title > pivot in the high part.
 *
 * @param bids Address of the vector<Bid> instance to be partitioned
 * @param begin Beginning index to partition (inclusive)
 * @param end   Ending index to partition (inclusive)
 * @return Index of the last element in the "low" partition (pivot position)
 */
int partition(vector<Bid>& bids, int begin, int end) {
    int low = begin;
    int high = end;

    // Pivot: use the title of the middle element to avoid worst-case O(n^2)
    // when the data is already sorted or nearly sorted.
    int mid = begin + (end - begin) / 2;
    string pivot = bids[mid].title;

    while (true) {
        // Move low right until we find an element that belongs in the high part
        while (low <= end && bids[low].title < pivot) {
            ++low;
        }
        // Move high left until we find an element that belongs in the low part
        while (high >= begin && pivot < bids[high].title) {
            --high;
        }
        // If pointers met or crossed, partitioning is done
        if (low >= high) {
            return high;
        }
        // Swap the misplaced elements and narrow the range
        swap(bids[low], bids[high]);
        ++low;
        --high;
    }
}

/**
 * Perform a quick sort on bid title (in-place).
 * Average performance: O(n log n)
 * Worst case performance: O(n^2) — e.g. already sorted with poor pivot choice
 *
 * @param bids  Address of the vector<Bid> instance to be sorted
 * @param begin Beginning index to sort (inclusive)
 * @param end   Ending index to sort (inclusive)
 */
void quickSort(vector<Bid>& bids, int begin, int end) {
    // Base case: 0 or 1 element is already sorted
    if (begin >= end) {
        return;
    }

    // Partition so that everything left of mid has title <= pivot,
    // everything right of mid has title >= pivot.
    int mid = partition(bids, begin, end);

    // Recursively sort the two partitions (exclude mid to avoid infinite loop)
    quickSort(bids, begin, mid);
    quickSort(bids, mid + 1, end);
}

//----------------------------------------------------------------------------
// FIXME (1a) — Selection sort
// Replace the stub body in VectorSorting.cpp (lines ~177–191).
//----------------------------------------------------------------------------

/**
 * Perform a selection sort on bid title (in-place).
 * Repeatedly finds the minimum title in the unsorted portion and swaps it
 * to the front. Average and worst case: O(n^2).
 *
 * @param bids Address of the vector<Bid> instance to be sorted
 */
void selectionSort(vector<Bid>& bids) {
    size_t size = bids.size();

    // Empty or single-element vector is already sorted
    if (size <= 1) {
        return;
    }

    // pos is the boundary: everything to the left is sorted
    for (size_t pos = 0; pos < size - 1; ++pos) {
        size_t min = pos;

        // Find the index of the smallest title in the unsorted part [pos .. size-1]
        for (size_t j = pos + 1; j < size; ++j) {
            if (bids[j].title < bids[min].title) {
                min = j;
            }
        }

        // Place the minimum at pos (swap with current element at pos)
        if (min != pos) {
            swap(bids[pos], bids[min]);
        }
    }
}

//----------------------------------------------------------------------------
// FIXME (1b) — Menu case 3: Selection sort and timing
// Replace the FIXME (1b) comment in main()'s switch (around line 265).
//----------------------------------------------------------------------------

        case 3:
            // Time the selection sort (O(n^2) — can be slow on large datasets)
            ticks = clock();
            selectionSort(bids);
            ticks = clock() - ticks;

            cout << bids.size() << " bids sorted with selection sort" << endl;
            cout << "time: " << ticks << " clock ticks" << endl;
            cout << "time: " << ticks * 1.0 / CLOCKS_PER_SEC << " seconds" << endl;
            break;

//----------------------------------------------------------------------------
// FIXME (2b) — Menu case 4: Quick sort and timing
// Replace the FIXME (2b) comment in main()'s switch (around line 267).
//----------------------------------------------------------------------------

        case 4:
            // Time the quick sort (O(n log n) average — faster on large datasets)
            ticks = clock();
            if (!bids.empty()) {
                quickSort(bids, 0, static_cast<int>(bids.size()) - 1);
            }
            ticks = clock() - ticks;

            cout << bids.size() << " bids sorted with quick sort" << endl;
            cout << "time: " << ticks << " clock ticks" << endl;
            cout << "time: " << ticks * 1.0 / CLOCKS_PER_SEC << " seconds" << endl;
            break;
