const gulp = require('gulp');
const sourcemaps = require('gulp-sourcemaps');
const babel = require('gulp-babel');
const concat = require('gulp-concat');
const spawn = require('child_process').spawn;

gulp.task('build', () =>
	gulp.src('src/**/*.js')
		.pipe(sourcemaps.init())
		.pipe(babel({
			presets: ['@babel/env']
		}))
		.pipe(concat('all.js'))
		.pipe(sourcemaps.write('.'))
		.pipe(gulp.dest('dist'))
);

gulp.task('serve', ['build'], () => {
	spawn('node', ['dist/all.js'], { stdio: 'inherit' });
});

gulp.task('default', ['build', 'serve']);