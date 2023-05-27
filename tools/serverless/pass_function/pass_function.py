def pass_function(event: dict, context: dict)->dict:
    finalGrade = context['finalGrade']
    
    return {"result": ", Congratulations on completing the course. Your score is {}!".format(finalGrade)}