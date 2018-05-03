package core.recordReaders;

import core.writables.ImageWritable;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FSDataOutputStream;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IOUtils;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.mapreduce.RecordWriter;
import org.apache.hadoop.mapreduce.TaskAttemptContext;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

import java.io.IOException;

/**
 * Created by Epanchee on 23.02.15.
 */
public abstract class ImageRecordWriter<I extends ImageWritable> extends RecordWriter<NullWritable, I> {

    protected final TaskAttemptContext taskAttemptContext;

    protected abstract void writeImage(I image, FSDataOutputStream imageFile) throws Exception;

    protected ImageRecordWriter(TaskAttemptContext taskAttemptContext){
        this.taskAttemptContext = taskAttemptContext;
    }

    @Override
    public void write(NullWritable nullWritable, I image) throws IOException, InterruptedException {
        if (image.getImage() != null) {
            FSDataOutputStream imageFile = null;
            Configuration job = taskAttemptContext.getConfiguration();
            Path file = FileOutputFormat.getOutputPath(taskAttemptContext);
            FileSystem fs = file.getFileSystem(job);
            // Constructing image filename and path
            Path imageFilePath = new Path(file, image.getFileName() + "."
                    + image.getFormat());

            try {
                // Creating file
                imageFile = fs.create(imageFilePath);
                writeImage(image, imageFile);

                // Write image to file using ImageIO
                //ImageIO.write(bufferedImageWritable.getBufferedImage(), bufferedImageWritable.getFormat(), imageFile);
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
                IOUtils.closeStream(imageFile);
            }
        }
    }

    @Override
    public void close(TaskAttemptContext taskAttemptContext) throws IOException, InterruptedException {

    }
}
