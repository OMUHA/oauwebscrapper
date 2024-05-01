-- Total students who passed -- 
-- Male Group -- 
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'M'
and phy_pts >  1 and chem_pts > 1 and bio_pts > 1 and bmath_pts > 1 and eng_pts > 1
 and   (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- Female ----------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'F'
  and phy_pts >  1 and chem_pts > 1 and bio_pts > 1 and bmath_pts > 1 and eng_pts > 1
  and   (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- Total students --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'
  and phy_pts >  1 and chem_pts > 1 and bio_pts > 1 and bmath_pts > 1 and eng_pts > 1
  and   (exam_year = '2022' or exam_year = '2023')
group by exam_year;


--- Students who failed --------------------------------
--Male group --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'M'
  and (phy_pts < 2  or chem_pts < 2  or  bio_pts < 2  or  bmath_pts < 2  or    eng_pts < 2  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;

-- Femail group --
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'F'
  and (phy_pts < 2  or chem_pts < 2  or  bio_pts < 2  or  bmath_pts < 2  or    eng_pts < 2  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- Total number of ----------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'
  and (phy_pts < 2  or chem_pts < 2  or  bio_pts < 2  or  bmath_pts < 2  or    eng_pts < 2  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;

-- Results summary for individual subjects --

-- Physics Male--
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'M'
  and (phy_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- Physics Female ---
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'F'
  and (phy_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- Phycis total --
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'
  and (phy_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;

--Chemistry Male --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'M'
  and (chem_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
--Chemistry Female --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'F'
  and (chem_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
--Chemistry Total --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'
  and (chem_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;

-- Biology Male --------------------------------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'M'
  and (bio_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
--Biology Female --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'F'
  and (bio_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
--Biology Total --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'
  and (bio_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;

-- Basic Mathematics Male-- --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'M'
  and (bmath_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- Basic Mathematics Female-- --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'F'
  and (bmath_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- Basic Mathematics Total-- --------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'
  and (bmath_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;

-- English lang Male--------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'M'
  and (eng_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- English lang Female--------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'  and candidate_gender = 'F'
  and (eng_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;
-- English lang Total--------------------------------
select exam_year,count(index_no) total from public.student_results
where  candidate_type = 'S'
  and (eng_pts >1  )
  and (exam_year = '2022' or exam_year = '2023')
group by exam_year;