def calculate_function(event: dict, context: dict)->dict:
    regularGrade = context['regularGrade']
    testGrade = context['testGrade']
    regularProportion = context['regularProportion']
    testProportion = context['testProportion']
    
    finalGrade = (regularGrade * regularProportion + testGrade * testProportion) / (regularProportion + testProportion)

    return {"finalGrade": finalGrade}    
