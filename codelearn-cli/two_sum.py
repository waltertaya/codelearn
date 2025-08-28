def two_sum(nums, target):
    """
    Given an array of integers nums and an integer target, 
    return indices of the two numbers such that they add up to target.
    """
    num_map = {}
    for i, num in enumerate(nums):
        complement = target - num
        if complement in num_map:
            return [num_map[complement], i]
        num_map[num] = i
    return []

# Test the function
if __name__ == "__main__":
    # Test case 1: [2,7,11,15], target = 9
    result1 = two_sum([2, 7, 11, 15], 9)
    print(f"Test 1: {result1}")  # Expected: [0, 1]
    
    # Test case 2: [3,2,4], target = 6
    result2 = two_sum([3, 2, 4], 6)
    print(f"Test 2: {result2}")  # Expected: [1, 2]
