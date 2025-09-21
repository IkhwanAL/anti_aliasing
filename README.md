# anti_aliasing
Experiment With Anti Aliasing

## What the motivation of this project
Well i overheard in youtube on how anti aliasing work, and as a person who playing games and always configure the graphic setting before playing.
The word Anti Aliasing is not a uncommon word, i always see that i know what it try to achive. To remove "jagged" line that look like staircase in 45 deg line or anything that not straigth line.
So as a programmer myself i'm curios on how does this Anti Aliasing work. So this is repo about experiment lab

## Journey
Ok, built an Anti Aliasing is an unknown world for me. But i do know the option e.g FXAA, MSAA TAA. But right now i will start simple that look similar to FXAA maybe below it. And the step to build are

1. Convert image to grayscale
2. Use That Grayscale Image to Create New Image That only Contain Edge (Edge Detection Image)
3. Use Edge Detection Image to Blend Edge By Using Blur (We Can Use Any Blur Here, But Most Recommended Are Gaussian Blur)

With This we create Anti Aliasing
