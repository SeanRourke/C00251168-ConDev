# Seán Rourke
# instructions for running
# navigate to folder containing file, run command python3 lab.py

def get_value(r):
    if r == 'I':
        return 1
    elif r == 'V':
        return 5
    elif r == 'X':
        return 10
    elif r == 'L':
        return 50
    elif r == 'C':
        return 100
    elif r == 'D':
        return 500
    elif r == 'M':
        return 1000
    else:
        return 0

def roman_to_int(n):
    if len(n) > 0 and len(n) < 16:
        int_val = 0
        i = 0
        valid = True
        while i < len(n):
            n1 = get_value(n[i])
            if n1 == 0:
                valid = False
            if i+1 < len(n):
                n2 = get_value(n[i+1])
                if n2 == 0:
                    valid = False
                if n1 >= n2:
                    int_val += n1
                else:
                    int_val += (n2-n1)
                    i +=1
            else:
                int_val += n1
            i += 1
        if valid == True:
            return int_val
        else:
            print("Invalid Input")
    else:
        print('Value out of range')

print(roman_to_int('III'))
print(roman_to_int('LVIII'))
print(roman_to_int('MCMXCIV'))