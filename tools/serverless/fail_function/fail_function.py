def fail_function(event: dict, context: dict)->dict:
    finalGrade = context['finalGrade']
    
    return {"result": "Unfortunately, you did not pass. Your score is {}!".format(finalGrade)}